package deployments

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	rolloutsv1alpha1 "github.com/nubank/hringhorni/pkg/apis/rollouts/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/logging"
)

const (
	ScopeControl    = "control"
	ScopeExperiment = "experiment"
	ScopeProduction = "production"

	RevisionAnnotationKey = "deployment.kubernetes.io/revision"
)

var (
	ErrReadinessDeadlineExceeded = errors.New("readiness deadline exceeded")
)

func NewBaseline(rollout *rolloutsv1alpha1.Rollout, productionDeployment *appsv1.Deployment, desiredReplicas int32) *appsv1.Deployment {
	baseline := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        kmeta.ChildName(rollout.Name, "-baseline"),
			Namespace:   rollout.Namespace,
			Labels:      make(map[string]string, len(productionDeployment.Labels)),
			Annotations: make(map[string]string, len(productionDeployment.Annotations)),
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(rollout),
			},
		},
		Spec: *productionDeployment.Spec.DeepCopy(),
	}
	baseline.Spec.Replicas = pointer.Int32(desiredReplicas)
	copyKeys(productionDeployment.Labels, baseline.Labels)
	copyKeys(productionDeployment.Annotations, baseline.Annotations)
	baseline.Labels[rolloutsv1alpha1.RolloutNameLabelKey] = rollout.Name
	if revision := rollout.Spec.Revision; revision != "" {
		baseline.Annotations[rolloutsv1alpha1.RolloutRevisionAnnotationKey] = revision
	}
	baseline.Labels[rolloutsv1alpha1.RolloutScopeLabelKey] = ScopeControl

	return baseline
}

func New(rollout *rolloutsv1alpha1.Rollout, desiredReplicas int32) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   rollout.Namespace,
			Labels:      make(map[string]string, len(rollout.Labels)+1),
			Annotations: make(map[string]string, len(rollout.Annotations)),
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(rollout),
			},
		},
		Spec: appsv1.DeploymentSpec{
			MinReadySeconds:         rollout.Spec.MinReadySeconds,
			ProgressDeadlineSeconds: pointer.Int32(int32(rollout.Spec.Timeout.Seconds())),
			Replicas:                pointer.Int32(desiredReplicas),
			Selector:                rollout.Spec.Selector,
			Template:                rollout.Spec.Template,
		},
	}
	copyKeys(rollout.Labels, deployment.Labels)
	copyKeys(rollout.Annotations, deployment.Annotations)
	deployment.Labels[rolloutsv1alpha1.RolloutNameLabelKey] = rollout.Name
	deployment.Labels[rolloutsv1alpha1.RolloutAppNameLabelKey] = rollout.Spec.AppName
	if revision := rollout.Spec.Revision; revision != "" {
		deployment.Annotations[rolloutsv1alpha1.RolloutRevisionAnnotationKey] = revision
	}
	return deployment
}

func copyKeys(in, out map[string]string) {
	for key, value := range in {
		out[key] = value
	}
}

func NewCanary(rollout *rolloutsv1alpha1.Rollout, desiredReplicas int32) *appsv1.Deployment {
	canary := New(rollout, desiredReplicas)
	canary.Name = kmeta.ChildName(rollout.Name, "-canary")
	canary.Labels[rolloutsv1alpha1.RolloutScopeLabelKey] = ScopeExperiment
	return canary
}

func NewProduction(rollout *rolloutsv1alpha1.Rollout, desiredReplicas int32) *appsv1.Deployment {
	production := New(rollout, desiredReplicas)
	production.Name = rollout.Name
	production.Labels[rolloutsv1alpha1.RolloutScopeLabelKey] = ScopeProduction
	production.OwnerReferences = nil
	return production
}

func Create(ctx context.Context, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	logger := logging.FromContext(ctx)
	kubeClient := kubeclient.Get(ctx)

	createdDeployment, err := kubeClient.AppsV1().Deployments(deployment.Namespace).Create(ctx, deployment, metav1.CreateOptions{})
	if err != nil {
		err = fmt.Errorf("error creating Deployment: %w", err)
		if apierrors.IsBadRequest(err) {
			err = controller.NewPermanentError(err)
		}
		return nil, err
	}

	logger.Info("Deployment has been successfully created")

	return createdDeployment, nil
}

func Update(ctx context.Context, deployment *appsv1.Deployment) error {
	logger := logging.FromContext(ctx)
	kubeClient := kubeclient.Get(ctx)

	_, err := kubeClient.AppsV1().Deployments(deployment.Namespace).Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		err = fmt.Errorf("error updating Deployment: %w", err)
		if apierrors.IsBadRequest(err) {
			err = controller.NewPermanentError(err)
		}
		return err
	}

	logger.Infof("Deployment %s has been successfully updated", deployment.Name)

	return nil
}

func Delete(ctx context.Context, deployment *appsv1.Deployment) error {
	logger := logging.FromContext(ctx)
	kubeClient := kubeclient.Get(ctx)

	err := kubeClient.AppsV1().Deployments(deployment.Namespace).Delete(ctx, deployment.Name, metav1.DeleteOptions{
		Preconditions: &metav1.Preconditions{
			UID:             &deployment.UID,
			ResourceVersion: &deployment.ResourceVersion,
		},
	})
	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Debugf("Deployment %s is no longer available in the namespace %s", deployment.Name, deployment.Namespace)
			return nil
		}
		err = fmt.Errorf("error deleting Deployment %s: %w", deployment.Name, err)
		if apierrors.IsBadRequest(err) {
			err = controller.NewPermanentError(err)
		}
		return err
	}

	logger.Infof("Deployment %s has been successfully deleted", deployment.Name)

	return nil
}

func Scale(ctx context.Context, deployment *appsv1.Deployment, replicas int32) error {
	logger := logging.FromContext(ctx)
	kubeClient := kubeclient.Get(ctx)

	patch := []struct {
		Op    string `json:"op"`
		Path  string `json:"path"`
		Value any    `json:"value"`
	}{{
		Op:    "replace",
		Path:  "/spec/replicas",
		Value: replicas,
	},
	}

	data, err := json.Marshal(patch)
	if err != nil {
		return controller.NewPermanentError(err)
	}

	logger.Infof("Scaling Deployment %s to %d replicas", deployment.Name, replicas)
	_, err = kubeClient.AppsV1().Deployments(deployment.Namespace).Patch(ctx, deployment.Name, types.JSONPatchType, data, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("error scaling Deployment %s to %d replicas: %w", deployment.Name, replicas, err)
	}
	return nil
}

func IsReady(deployment *appsv1.Deployment, desiredReplicas int32) (bool, error) {
	if deployment.Generation != deployment.Status.ObservedGeneration {
		return false, nil
	}
	for _, condition := range deployment.Status.Conditions {
		if condition.Type == appsv1.DeploymentProgressing && condition.Status == corev1.ConditionFalse && condition.Reason == "ProgressDeadlineExceeded" {
			return false, fmt.Errorf("%w: Deployment %s did not become ready in %d seconds", ErrReadinessDeadlineExceeded, deployment.Name, *deployment.Spec.ProgressDeadlineSeconds)
		}

		if condition.Type == appsv1.DeploymentAvailable {
			if condition.Status == corev1.ConditionFalse && time.Since(condition.LastUpdateTime.Time).Seconds() >= float64(*deployment.Spec.ProgressDeadlineSeconds) {
				return false, fmt.Errorf("%w: Deployment %s did not become ready in %d seconds", ErrReadinessDeadlineExceeded, deployment.Name, *deployment.Spec.ProgressDeadlineSeconds)
			}

			if condition.Status == corev1.ConditionTrue && deployment.Status.AvailableReplicas == desiredReplicas {
				return true, nil
			}
		}
	}
	return false, nil
}
