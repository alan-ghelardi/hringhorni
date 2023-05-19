package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Rollout struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RolloutSpec   `json:"spec,omitempty"`
	Status            RolloutStatus `json:"status,omitempty"`
}

var (
	// Check that Rollout can be validated and defaulted.
	_ kmeta.OwnerRefable = (*Rollout)(nil)
	// Check that the type conforms to the duck Knative Resource shape.
	_ duckv1.KRShaped = (*Rollout)(nil)
)

// RolloutSpec describes the desired state of this Rollout.
type RolloutSpec struct {
	AppName string `json:"appName"`

	// +optional
	Revision string `json:"revision,omitempty"`

	// +optional
	ExternalLink string `json:"externalLink,omitempty"`

	// +optional
	Canary *CanarySettings `json:"canary,omitempty"`

	// +optional
	Actions RolloutActions `json:"actions,omitempty"`

	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	Template corev1.PodTemplateSpec `json:"template"`

	// Minimum number of seconds for which a newly created pod should be ready
	// without any of its container crashing, for it to be considered available.
	// Defaults to 0 (pod will be considered available as soon as it is ready)
	// +optional
	MinReadySeconds int32 `json:"minReadySeconds,omitempty"`

	// +optional
	Timeout *metav1.Duration `json:"timeout,omitempty"`
}

type CanarySettings struct {
	Percentage *float64 `json:"percentage,omitempty"`

	Duration *metav1.Duration `json:"duration,omitempty"`

	Interval *metav1.Duration `json:"interval,omitempty"`
}

// RolloutActions represents user requested actions that can be performed during
// a rollout.
type RolloutActions struct {

	// Cancels an ongoing rollout.
	// +optional
	Cancel bool `json:"cancel"`

	// Rolls back a completed rollout.
	// +optional
	RollBack bool `json:"rollBack"`
}

// RolloutStatus describes the current state of this Rollout.
type RolloutStatus struct {
	duckv1.Status `json:",inline"`

	// Relevant information about the previous rollout of the application in
	// question.
	// +optional
	PreviousRollout *PreviousRollout `json:"previousRollout,omitempty"`

	// Timestamp at which the rollout has started.
	// +optional
	StartedAt *metav1.Time `json:"startedAt,omitempty"`

	// Timestamp at which the rollout has completed.
	// +optional
	CompletedAt *metav1.Time `json:"completedAt,omitempty"`

	// Status of the canary rollout.
	// +optional
	Canary *CanaryStatus `json:"canary,omitempty"`

	// Whether or not the controller has started rolling back the underlying
	// Deployment.
	RollBackRequested bool `json:"rollBackRequested"`
}

type PreviousRollout struct {
	DeploymentName string `json:"deploymentName"`

	DeploymentRevision string `json:"deploymentRevision,omitempty"`

	Revision string `json:"revision,omitempty"`

	Replicas int32 `json:"replicas"`
}

type CanaryStatus struct {
	BaselineDeployment *string `json:"baselineDeployment,omitempty"`

	CanaryDeployment *string `json:"canaryDeployment,omitempty"`

	Completed bool `json:"completed"`
}

func (c *CanaryStatus) Iscompleted() bool {
	return c != nil && c.Completed
}

// GetStatus retrieves the status of the resource. Implements the KRShaped interface.
func (r *Rollout) GetStatus() *duckv1.Status {
	return &r.Status.Status
}

func (r *Rollout) IsDone() bool {
	return !r.Status.GetCondition(RolloutConditionSucceeded).IsUnknown()
}

func (r *Rollout) IsRollingBack() bool {
	condition := r.Status.GetCondition(RollbackConditionSucceeded)
	if condition == nil && r.Spec.Actions.RollBack {
		return true
	}
	return condition.IsUnknown() && condition.GetReason() != ""
}

// RolloutList is a list of Rollout objects.
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type RolloutList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Rollout `json:"items"`
}
