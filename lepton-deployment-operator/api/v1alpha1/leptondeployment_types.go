/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"fmt"

	"github.com/leptonai/lepton/go-pkg/util"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// LeptonDeploymentSpec defines the desired state of LeptonDeployment
type LeptonDeploymentSpec struct {
	LeptonDeploymentSystemSpec `json:",inline"`
	LeptonDeploymentUserSpec   `json:",inline"`
}

// LeptonDeploymentStatus defines the system-controlled spec.
type LeptonDeploymentSystemSpec struct {
	PhotonName  string `json:"photon_name"`
	PhotonImage string `json:"photon_image"`
	BucketName  string `json:"bucket_name"`
	// +optional
	EFSID string `json:"efs_id"`
	// +optional
	EFSAccessPointID string `json:"efs_access_point_id"`
	PhotonPrefix     string `json:"photon_prefix"`
	// TODO: delete ServiceAccountName at some point, which is only used for backward compatibility.
	ServiceAccountName            string `json:"service_account_name,omitempty"`
	S3ReadOnlyAccessK8sSecretName string `json:"s3_read_only_access_k8s_secret_name,omitempty"`
	RootDomain                    string `json:"root_domain,omitempty"`
	WorkspaceName                 string `json:"workspace_name,omitempty"`
	WorkspaceToken                string `json:"workspace_token,omitempty"`
	CertificateARN                string `json:"certificate_arn,omitempty"`
}

// LeptonDeploymentStatus defines the user-controlled spec.
type LeptonDeploymentUserSpec struct {
	Name                string                              `json:"name"`
	PhotonID            string                              `json:"photon_id"`
	ResourceRequirement LeptonDeploymentResourceRequirement `json:"resource_requirement"`
	// +optional
	APITokens []TokenVar `json:"api_tokens"`
	// +optional
	Envs []EnvVar `json:"envs"`
	// +optional
	Mounts []Mount `json:"mounts"`
}

// GetSpecName returns the name of the deployment.
func (ld LeptonDeployment) GetSpecName() string {
	return ld.Spec.Name
}

// GetSpecID returns the ID of the deployment. It equals to the Name.
func (ld LeptonDeployment) GetSpecID() string {
	return ld.GetSpecName()
}

// GetVersion returns the version of the deployment, which is always 0 because we don't support versioning.
func (ld LeptonDeployment) GetVersion() int64 {
	return 0
}

// Get Shape returns the resource shape of the deployment as a string.
func (ld *LeptonDeployment) GetShape() string {
	if ld.Spec.ResourceRequirement.ResourceShape == "" {
		return string(Customized)
	}
	return string(ld.Spec.ResourceRequirement.ResourceShape)
}

func (ld *LeptonDeployment) GetTokens() []string {
	tokens := make([]string, 0)
	for _, tokenVar := range ld.Spec.APITokens {
		if tokenVar.Value != "" {
			tokens = append(tokens, tokenVar.Value)
		} else if tokenVar.ValueFrom.TokenNameRef == TokenNameRefWorkspaceToken {
			if ld.Spec.WorkspaceToken != "" {
				tokens = append(tokens, ld.Spec.WorkspaceToken)
			}
		}
	}
	return util.UniqStringSlice(tokens)
}

// Patch modifies the deployment with the given user spec. It only supports PhotonID and MinReplicas for now.
func (ld *LeptonDeployment) Patch(p *LeptonDeploymentUserSpec) {
	if p.PhotonID != "" {
		ld.Spec.PhotonID = p.PhotonID
	}
	if p.ResourceRequirement.MinReplicas > 0 {
		ld.Spec.ResourceRequirement.MinReplicas = p.ResourceRequirement.MinReplicas
	}
}

// LeptonDeploymentResourceRequirement defines the resource requirement of the deployment.
type LeptonDeploymentResourceRequirement struct {
	LeptonDeploymentReplicaResourceRequirement `json:",inline"`
	// +optional
	ResourceShape LeptonDeploymentResourceShape `json:"resource_shape"`
	MinReplicas   int32                         `json:"min_replicas"`
}

// GetAcceleratorRequirement returns the required number and type of the accelerator for the deployment.
func (lr *LeptonDeploymentResourceRequirement) GetAcceleratorRequirement() (float64, string) {
	if lr.ResourceShape != "" {
		r, err := ShapeToReplicaResourceRequirement(lr.ResourceShape)
		if err != nil {
			return 0, ""
		}
		return r.AcceleratorNum, r.AcceleratorType
	}

	return lr.AcceleratorNum, lr.AcceleratorType
}

// LeptonDeploymentReplicaResourceRequirement defines the resource requirement of the deployment.
type LeptonDeploymentReplicaResourceRequirement struct {
	CPU    float64 `json:"cpu"`
	Memory int64   `json:"memory"`
	// +optional
	AcceleratorType string `json:"accelerator_type"`
	// +optional
	AcceleratorNum float64 `json:"accelerator_num"`
	// +optional
	EphemeralStorageInGB int64 `json:"ephemeral_storage_in_gb"`
}

// ResourceShape defines the resource shape of the deployment.
type ResourceShape struct {
	// Name of the shape. E.g. "Large"
	Name string `json:"name"`
	// Description of the shape. E.g. "Large shape with 4 CPUs and 16GB of RAM"
	Description string                                     `json:"description"`
	Resource    LeptonDeploymentReplicaResourceRequirement `json:"resource"`
}

// TokenVar defines the token variable of the deployment.
type TokenVar struct {
	Value     string     `json:"value,omitempty"`
	ValueFrom TokenValue `json:"value_from,omitempty"`
}

// TokenValue defines the token value of the deployment.
type TokenValue struct {
	TokenNameRef string `json:"token_name_ref,omitempty"`
	// TODO: we can add SecretNameRef if we want to support reading token from secret.
}

const (
	// TokenNameRefWorkspaceToken is the token name ref for workspace token.
	// It refers to the current workspace token. Updating the workspace token
	// will not refresh the deployment token.
	TokenNameRefWorkspaceToken = "WORKSPACE_TOKEN"
)

// EnvVar defines the environment variable of the deployment.
type EnvVar struct {
	Name      string   `json:"name"`
	Value     string   `json:"value,omitempty"`
	ValueFrom EnvValue `json:"value_from,omitempty"`
}

// EnvValue defines the environment variable value of the deployment.
type EnvValue struct {
	SecretNameRef string `json:"secret_name_ref,omitempty"`
}

// Mount defines the volumes of the deployment.
type Mount struct {
	Path      string `json:"path"`
	MountPath string `json:"mount_path"`
}

// LeptonDeploymentStatus defines the observed state of LeptonDeployment
type LeptonDeploymentStatus struct {
	State    LeptonDeploymentState    `json:"state"`
	Endpoint LeptonDeploymentEndpoint `json:"endpoint"`
}

// LeptonDeploymentResourceShape defines the resource shape of the deployment.
type LeptonDeploymentResourceShape string

// LeptonDeploymentState defines the state of the deployment.
type LeptonDeploymentState string

const (
	LeptonDeploymentStateRunning  LeptonDeploymentState = "Running"
	LeptonDeploymentStateNotReady LeptonDeploymentState = "Not Ready"
	LeptonDeploymentStateStarting LeptonDeploymentState = "Starting"
	LeptonDeploymentStateUpdating LeptonDeploymentState = "Updating"
	LeptonDeploymentStateDeleting LeptonDeploymentState = "Deleting"
	LeptonDeploymentStateUnknown  LeptonDeploymentState = ""
)

// DeletionFinalizerName is the name of the finalizer for deletion. It is used to
// prevent the deletion of the LeptonDeployment before cooresponding resources are
// deleted, including deployment, service, and ingress.
const DeletionFinalizerName = "lepton.ai/deletion"

// LeptonDeploymentEndpoint defines the endpoint of the deployment.
type LeptonDeploymentEndpoint struct {
	InternalEndpoint string `json:"internal_endpoint"`
	ExternalEndpoint string `json:"external_endpoint"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=ld

// LeptonDeployment is the Schema for the leptondeployments API
type LeptonDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LeptonDeploymentSpec   `json:"spec,omitempty"`
	Status LeptonDeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LeptonDeploymentList contains a list of LeptonDeployment
type LeptonDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LeptonDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LeptonDeployment{}, &LeptonDeploymentList{})
}

// ShapeToReplicaResourceRequirement converts a shape to a replica resource requirement.
func ShapeToReplicaResourceRequirement(shape LeptonDeploymentResourceShape) (*LeptonDeploymentReplicaResourceRequirement, error) {
	shape = DisplayShapeToShape(string(shape))
	s := SupportedShapesAWS[shape]
	if s == nil {
		return nil, fmt.Errorf("shape %s is not supported", shape)
	}

	return &s.Resource, nil
}
