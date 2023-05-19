package v1alpha1

import (
	"context"

	"knative.dev/pkg/apis"
)

// Validate implements apis.Validatable.
func (a *Analysis) Validate(ctx context.Context) (errors *apis.FieldError) {
	// We don't want to block the object deletion if it's in an invalid
	// state
	if apis.IsInDelete(ctx) {
		return nil
	}
	return apis.ValidateObjectMetadata(&a.ObjectMeta).ViaField("metadata").
		Also(a.Spec.Validate(apis.WithinSpec(ctx)).ViaField("spec"))
}

// Validate implements apis.Validatable
func (as *AnalysisSpec) Validate(ctx context.Context) (errors *apis.FieldError) {
	return
}
