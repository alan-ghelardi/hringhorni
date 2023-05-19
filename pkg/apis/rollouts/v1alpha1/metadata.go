package v1alpha1

// Common labels and annotations used across the codebase.
var (
	RolloutNameLabelKey          = withGroup("rollout-name")
	RolloutAppNameLabelKey       = withGroup("app-name")
	RolloutScopeLabelKey         = withGroup("scope")
	RolloutRevisionAnnotationKey = withGroup("revision")
)

func withGroup(name string) string {
	return SchemeGroupVersion.Group + "/" + name
}
