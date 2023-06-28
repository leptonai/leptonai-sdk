package gcssyncer

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/leptonai/lepton/go-pkg/k8s"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateSyncerForDefaultEFS creates a deployment that syncs the given GCS URL to the given path under the default EFS root.
func CreateSyncerForDefaultEFS(ns, name string, gcsURL, path, credJSON string) error {
	image := "google/cloud-sdk"
	mountPath := "/mnt/efs/default/"
	volumeName := "default-efs"
	claimName := ns + "-efs-default-pvc"
	encoded := base64.StdEncoding.EncodeToString([]byte(credJSON))

	command := fmt.Sprintf("credFilePath=$(mktemp) && echo '%s' | base64 --decode > \"$credFilePath\" && gcloud auth login --cred-file=$credFilePath && while true; do gsutil -m rsync -r %s %s; sleep 1; done",
		encoded, gcsURL, "/mnt/efs/default"+path)
	appName := "gcs-syncer"

	labels := map[string]string{
		"app":    appName,
		"syncer": name,
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
	}

	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
		},
		Spec: corev1.PodSpec{
			SecurityContext: k8s.DefaultPodSecurityContext(),
			Containers: []corev1.Container{
				{
					Name:  appName,
					Image: image,
					Command: []string{
						"/bin/sh",
						"-c",
						command,
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      volumeName,
							MountPath: mountPath,
						},
					},
					SecurityContext: k8s.DefaultContainerSecurityContext(),
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: volumeName,
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: claimName,
						},
					},
				},
			},
		},
	}

	deployment.Spec.Template = podTemplate

	return k8s.Client.Create(context.TODO(), deployment)
}

// DeleteSyncerForDefaultEFS deletes the deployment with the given name in the given namespace.
func DeleteSyncerForDefaultEFS(ns, name string) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
	}

	return k8s.Client.Delete(context.Background(), deployment)
}
