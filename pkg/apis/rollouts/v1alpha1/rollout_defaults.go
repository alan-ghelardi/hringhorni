package v1alpha1

import (
	"context"

	"k8s.io/utils/pointer"
)

// SetDefaults implements apis.Defaultable.
func (r *Rollout) SetDefaults(ctx context.Context) {
	r.Spec.SetDefaults(ctx)
}

// SetDefaults implements apis.Defaultable.
func (rs *RolloutSpec) SetDefaults(ctx context.Context) {
	if rs.Replicas == nil {
		rs.Replicas = pointer.Int32(1)
	}
}
