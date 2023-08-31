package httpapi

import (
	"context"
	"fmt"
	"time"

	leptonaiv1alpha1 "github.com/leptonai/lepton/deployment-operator/api/v1alpha1"
	"github.com/leptonai/lepton/go-pkg/datastore"
	"github.com/leptonai/lepton/go-pkg/k8s/secret"
	goutil "github.com/leptonai/lepton/go-pkg/util"

	"gocloud.dev/blob"
	"k8s.io/utils/ptr"
)

// Backupable is an interface for resources that can be backed up
type Backupable interface {
	Backup(ctx context.Context) error
}

type Handler struct {
	clusterName                   string
	namespace                     string
	prometheusURL                 string
	bucketName                    string
	efsID                         string
	photonPrefix                  string
	photonImageRegistry           string
	s3ReadOnlyAccessK8sSecretName string
	rootDomain                    string
	sharedALBMainDomain           string
	workspaceName                 string
	certARN                       string
	apiToken                      string
	photonBucket                  *blob.Bucket

	secretDB *secret.SecretSet
	phDB     *datastore.CRStore[*leptonaiv1alpha1.Photon]
	ldDB     *datastore.CRStore[*leptonaiv1alpha1.LeptonDeployment]
}

func New(clusterName, namespace, prometheusURL, bucketName, efsID, protonPrefix, photonImageRegistry, s3ReadOnlyAccessK8sSecretName,
	rootDomain, sharedALBMainDomain, workspaceName, certARN, apiToken string, photonBucket, backupBucket *blob.Bucket,
	workspaceState WorkspaceState) *Handler {
	h := &Handler{
		clusterName:                   clusterName,
		namespace:                     namespace,
		prometheusURL:                 prometheusURL,
		bucketName:                    bucketName,
		efsID:                         efsID,
		photonPrefix:                  protonPrefix,
		photonImageRegistry:           photonImageRegistry,
		s3ReadOnlyAccessK8sSecretName: s3ReadOnlyAccessK8sSecretName,
		rootDomain:                    rootDomain,
		sharedALBMainDomain:           sharedALBMainDomain,
		workspaceName:                 workspaceName,
		certARN:                       certARN,
		apiToken:                      apiToken,
		photonBucket:                  photonBucket,

		secretDB: secret.New(namespace, secret.SecretObjectName, backupBucket),
		phDB: datastore.NewCRStore[*leptonaiv1alpha1.Photon](
			namespace,
			&leptonaiv1alpha1.Photon{},
			backupBucket,
		),
		ldDB: datastore.NewCRStore[*leptonaiv1alpha1.LeptonDeployment](
			namespace,
			&leptonaiv1alpha1.LeptonDeployment{},
			backupBucket,
		),
	}

	switch workspaceState {
	case WorkspaceStateNormal:
	case WorkspaceStatePaused:
	case WorkspaceStateTerminated:
		go h.scaleAllDeploymentsDownToZeroReplica()
	default:
		goutil.Logger.Fatalw("invalid workspace state",
			"state", workspaceState,
		)
	}
	goutil.Logger.Infow("workspace state",
		"state", workspaceState,
	)

	go runBackupsPeriodically(h.phDB)
	go runBackupsPeriodically(h.ldDB)
	go runBackupsPeriodically(h.secretDB)

	return h
}

const (
	scaleAllDeploymentsDownToZeroReplicaRetryInterval = 5 * time.Minute
	scaleAllDeploymentsDownToZeroReplicaIdleInterval  = 24 * time.Hour
)

func (h *Handler) scaleAllDeploymentsDownToZeroReplica() {
	// first, backup the existing deployments
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		err := h.ldDB.Backup(ctx)
		cancel()
		if err == nil {
			break
		}
		goutil.Logger.Errorw(fmt.Sprintf("failed to backup, retrying in %v", backupRetryInterval),
			"namespace", h.namespace,
			"operation", "scaleAllDeploymentsDownToZeroReplica",
			"error", err,
		)
		time.Sleep(scaleAllDeploymentsDownToZeroReplicaRetryInterval)
	}
	goutil.Logger.Infow("backup successful",
		"operation", "scaleAllDeploymentsDownToZeroReplica",
		"namespace", h.namespace,
	)
	// second, scale all deployments down to zero replica
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		lds, err := h.ldDB.List(ctx)
		if err != nil {
			goutil.Logger.Errorw(fmt.Sprintf("failed to list deployments, retrying in %v", backupRetryInterval),
				"namespace", h.namespace,
				"operation", "scaleAllDeploymentsDownToZeroReplica",
				"error", err,
			)
			cancel()
			time.Sleep(scaleAllDeploymentsDownToZeroReplicaRetryInterval)
			continue
		}
		failed := false
		for _, ld := range lds {
			if ld.Spec.ResourceRequirement.MinReplicas == nil {
				continue
			}
			if *ld.Spec.ResourceRequirement.MinReplicas == 0 {
				continue
			}
			ld.Spec.ResourceRequirement.MinReplicas = ptr.To[int32](0)
			err = h.ldDB.Update(ctx, ld.GetSpecName(), ld)
			if err != nil {
				goutil.Logger.Errorw(fmt.Sprintf("failed to update deployment, retrying in %v", backupRetryInterval),
					"operation", "scaleAllDeploymentsDownToZeroReplica",
					"namespace", h.namespace,
					"deployment", ld.GetSpecName(),
					"error", err,
				)
				failed = true
				continue
			}
			goutil.Logger.Infow("scaled deployment down to zero replica",
				"operation", "scaleAllDeploymentsDownToZeroReplica",
				"namespace", h.namespace,
				"deployment", ld.GetSpecName(),
			)
		}
		cancel()
		if failed {
			time.Sleep(scaleAllDeploymentsDownToZeroReplicaRetryInterval)
		} else {
			time.Sleep(scaleAllDeploymentsDownToZeroReplicaIdleInterval)
		}
	}
}

func (h *Handler) PhotonHanlder() *PhotonHandler {
	return &PhotonHandler{
		Handler: *h,
	}
}

func (h *Handler) DeploymentHandler() *DeploymentHandler {
	return &DeploymentHandler{
		Handler: *h,
	}
}

func (h *Handler) MonitoringHandler() *MonitorningHandler {
	return &MonitorningHandler{
		Handler: *h,
	}
}

func (h *Handler) ReplicaHandler() *ReplicaHandler {
	return &ReplicaHandler{
		Handler: *h,
	}
}

func (h *Handler) SecretHandler() *SecretHandler {
	return &SecretHandler{
		Handler: *h,
	}
}

func (h *Handler) DeploymentEventHandler() *DeploymentEventHandler {
	return &DeploymentEventHandler{
		Handler: *h,
	}
}

func (h *Handler) StorageHandler() *StorageHandler {
	return &StorageHandler{
		Handler: *h,
	}
}

func (h *Handler) StorageSyncerHandler() *StorageSyncerHandler {
	return &StorageSyncerHandler{
		Handler: *h,
	}
}

func (h *Handler) DeploymentReadinessHandler() *DeploymentReadinessHandler {
	return &DeploymentReadinessHandler{
		Handler: *h,
	}
}

func (h *Handler) DeploymentTerminationHandler() *DeploymentTerminationHandler {
	return &DeploymentTerminationHandler{
		Handler: *h,
	}
}

const (
	backupPeriod        = 12 * time.Hour
	backupRetryInterval = 5 * time.Minute
)

func runBackupsPeriodically(c Backupable) {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		err := c.Backup(ctx)
		cancel()
		if err != nil {
			goutil.Logger.Errorw(fmt.Sprintf("failed to backup, retrying in %v", backupRetryInterval),
				"operation", "backup",
				"error", err,
			)
			time.Sleep(backupRetryInterval)
			continue
		}

		goutil.Logger.Infow(fmt.Sprintf("backup successful, next backup after %v", backupPeriod),
			"operation", "backup",
		)
		time.Sleep(backupPeriod)
	}
}
