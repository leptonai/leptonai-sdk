package main

import (
	"context"
	"net/http"
	"time"

	"github.com/leptonai/lepton/go-pkg/httperrors"
	"github.com/leptonai/lepton/go-pkg/k8s"
	"github.com/leptonai/lepton/go-pkg/util"
	goutil "github.com/leptonai/lepton/go-pkg/util"
	"github.com/leptonai/lepton/lepton-mothership/cluster"
	"github.com/leptonai/lepton/lepton-mothership/httpapi"
	"github.com/leptonai/lepton/lepton-mothership/terraform"
	"github.com/leptonai/lepton/lepton-mothership/workspace"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "gocloud.dev/blob/s3blob"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
)

func main() {
	terraform.MustInit()

	cluster.Init()
	workspace.Init()

	router := gin.Default()

	logger := goutil.Logger.Desugar()
	// Add a ginzap middleware, which:
	//   - Logs all requests, like a combined access and error log.
	//   - Logs to stdout.
	//   - RFC3339 with UTC time format.
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	// Logs all panic to error log
	//   - stack means whether output the stack info.
	router.Use(ginzap.RecoveryWithZap(logger, true))

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := router.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/info", httpapi.NewInfoHandler().HandleGet)

	v1.GET("/clusters", httpapi.HandleClusterList)
	v1.POST("/clusters", httpapi.HandleClusterCreate)
	v1.GET("/clusters/:clname", httpapi.HandleClusterGet)
	v1.GET("/clusters/:clname/logs", httpapi.HandleClusterGetLogs)
	v1.GET("/clusters/:clname/failure", httpapi.HandleClusterGetFailureLog)
	v1.PATCH("/clusters", httpapi.HandleClusterUpdate)
	v1.DELETE("/clusters/:clname", httpapi.HandleClusterDelete)

	v1.GET("/workspaces", httpapi.HandleWorkspaceList)
	v1.POST("/workspaces", httpapi.HandleWorkspaceCreate)
	v1.GET("/workspaces/:wsname", httpapi.HandleWorkspaceGet)
	v1.GET("/workspaces/:wsname/logs", httpapi.HandleWorkspaceGetLogs)
	v1.GET("//workspaces/:wsname/failure", httpapi.HandleWorkspaceGetFailureLog)
	v1.PATCH("/workspaces", httpapi.HandleWorkspaceUpdate)
	v1.DELETE("/workspaces/:wsname", httpapi.HandleWorkspaceDelete)

	v1.GET("/users", httpapi.HandleUserList)
	v1.POST("/users", httpapi.HandleUserCreate)
	v1.GET("/users/:uname", httpapi.HandleUserGet)
	v1.DELETE("/users/:uname", httpapi.HandleUserDelete)

	v1.PUT("/upgrade/:imagetag", upgradeHandler)

	if err := router.Run(":15213"); err != nil {
		goutil.Logger.Fatalw("Failed to start mothership",
			"operation", "router.Run",
			"error", err,
		)
	}
}

func upgradeHandler(c *gin.Context) {
	imageTag := c.Param("imagetag")
	if imageTag == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    httperrors.ErrorCodeInvalidRequest,
			"message": "Image tag is required",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Requested to restart. It may take some time.",
	})
	// start it in a new go routine to allow the http request to finish before shutting down
	go waitForIdlingAndUpdateImageTag(imageTag)
}

func waitForIdlingAndUpdateImageTag(imageTag string) {
	goutil.Logger.Infow("Received request to upgrade mothership",
		"operation", "updateImageTag",
		"image_tag", imageTag,
	)
	// wait until no jobs are running
	for {
		cluster.Worker.Lock()
		workspace.Worker.Lock()
		clusterJobsCount := cluster.Worker.CountJobs()
		workspaceJobsCount := workspace.Worker.CountJobs()
		goutil.Logger.Infow("Checking if any jobs are running",
			"operation", "updateImageTag",
			"cluster_jobs_count", clusterJobsCount,
			"workspace_jobs_count", workspaceJobsCount,
		)
		if clusterJobsCount == 0 && workspaceJobsCount == 0 {
			err := updateImageTag(imageTag)
			if err != nil {
				goutil.Logger.Errorw("Failed to update mothership deployment image tag, will retry",
					"operation", "updateImageTag",
					"error", err,
				)
			} else {
				goutil.Logger.Infow("Updated the image of mothership deployment, exiting...",
					"operation", "updateImageTag",
					"image_tag", imageTag,
				)
				// Done with the update, waiting for k8s to bring up a new pod and terminate this one.
				// Do not release locks to prevent new jobs from being created.
				// Requests sent during this time will be scheduled after restart.

				// TODO: requests coming in during this time will be queued but the http request will not finish.
				// Need to find a way to return a response.
				return
			}
		}
		workspace.Worker.Unlock()
		cluster.Worker.Unlock()
		time.Sleep(time.Minute)
	}
}

func updateImageTag(imageTag string) error {
	deployment := appsv1.Deployment{}
	// TODO: remove the hardcoded name and namespace
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err := k8s.Client.Get(ctx, types.NamespacedName{Name: "mothership", Namespace: "default"}, &deployment)
	if err != nil {
		return err
	}
	// update the image tag of the deployment
	image := deployment.Spec.Template.Spec.Containers[0].Image
	deployment.Spec.Template.Spec.Containers[0].Image = util.UpdateImageTag(image, imageTag)
	return k8s.Client.Update(ctx, &deployment)
}
