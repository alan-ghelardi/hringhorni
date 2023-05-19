package rollout

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/kmeta"
	"knative.dev/pkg/logging"

	hringhorniclientset "github.com/nubank/hringhorni/pkg/client/clientset/versioned"
	rolloutlisters "github.com/nubank/hringhorni/pkg/client/listers/rollouts/v1alpha1"
	"github.com/nubank/hringhorni/pkg/rolloutstrategy"

	k8sclock "k8s.io/utils/clock"
	"k8s.io/utils/pointer"

	"k8s.io/client-go/kubernetes"

	appslisters "k8s.io/client-go/listers/apps/v1"

	"github.com/nubank/hringhorni/pkg/reconciler/events"

	rolloutsv1alpha1 "github.com/nubank/hringhorni/pkg/apis/rollouts/v1alpha1"
	rolloutreconciler "github.com/nubank/hringhorni/pkg/client/injection/reconciler/rollouts/v1alpha1/rollout"
	"github.com/nubank/hringhorni/pkg/deployments"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/reconciler"
)

const (
	defaultRequeueDuration = 15 * time.Second
)

var (
	// clock allows us to control the clock for testing purposes.
	clock k8sclock.PassiveClock = k8sclock.RealClock{}
)

// Reconciler implements rolloutreconciler.Interface for
// Rollout resources.
type Reconciler struct {

	// analysisLister allows us to read Analysis objects from the indexer.
	analysisLister rolloutlisters.AnalysisLister

	// deploymentLister allows us to read Deployments from the indexer.
	deploymentLister appslisters.DeploymentLister

	// hringhorniClient allows us to talk to the Kubernetes API server for
	// Hringhorni resources.
	hringhorniClient hringhorniclientset.Interface

	// kubeClient allows us to talk to the Kubernetes API server for core
	// resources.
	kubeClient kubernetes.Interface

	// canary allows us to run a canary rollout.
	canary *rolloutstrategy.Canary
}

// Check that our Reconciler implements the required interfaces
var (
	_ rolloutreconciler.Interface = (*Reconciler)(nil)
)

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, rollout *rolloutsv1alpha1.Rollout) reconciler.Event {
	logger := logging.FromContext(ctx)
	logger.Info("Reconciling Rollout")

	ctx = context.WithValue(ctx, kubeclient.Key{}, r.kubeClient)

	if rollout.IsDone() && !rollout.IsRollingBack() {
		logger.Info("Rollout is done - no further actions to be performed")
		return controller.NewSkipKey(rollout.Namespace + "/" + rollout.Name)
	}

	previousCondition := rollout.Status.GetCondition(rolloutsv1alpha1.RolloutConditionSucceeded)

	err := r.reconcileRollout(ctx, rollout)
	if err != nil {
		handleError(ctx, rollout, err)
	}

	events.Emit(ctx, previousCondition, rollout.Status.GetCondition(rolloutsv1alpha1.RolloutConditionSucceeded), rollout)

	return err
}

func handleError(ctx context.Context, rollout *rolloutsv1alpha1.Rollout, err error) {
	logger := logging.FromContext(ctx)
	logger.Debugw("The controller encountered an error reconciling the Rollout", zap.Error(err))
	if controller.IsPermanentError(err) {
		if rollout.IsDone() {
			rollout.Status.MarkRollbackFailed(rolloutsv1alpha1.RollbackErrorReason, err.Error())
		} else {
			rollout.Status.MarkFailed(rolloutsv1alpha1.RolloutErrorReason, err.Error())
		}
	} else if isRequeueKey, _ := controller.IsRequeueKey(err); !isRequeueKey {
		if rollout.IsDone() {
			rollout.Status.MarkRollbackUnknown(rolloutsv1alpha1.RollbackErrorReason, err.Error())
		} else {
			rollout.Status.MarkUnknown(rolloutsv1alpha1.RolloutErrorReason, err.Error())
		}
	}
}

func (r *Reconciler) reconcileRollout(ctx context.Context, rollout *rolloutsv1alpha1.Rollout) error {
	logger := logging.FromContext(ctx)

	productionDeployment, err := r.getProductionDeployment(rollout)
	if err != nil {
		return err
	}

	// Store information about the previous rollout if it's still unset.
	if productionDeployment != nil && rollout.Status.PreviousRollout == nil {
		rollout.Status.PreviousRollout = &rolloutsv1alpha1.PreviousRollout{
			DeploymentName:     productionDeployment.Name,
			DeploymentRevision: productionDeployment.Annotations[deployments.RevisionAnnotationKey],
			Replicas:           *productionDeployment.Spec.Replicas,
			Revision:           productionDeployment.Annotations[rolloutsv1alpha1.RolloutRevisionAnnotationKey],
		}
	}

	if rollout.IsRollingBack() {
		logger.Debug("Continuing rollback process")
		err := r.rollBack(ctx, rollout, productionDeployment)
		if errors.Is(err, deployments.ErrReadinessDeadlineExceeded) {
			rollout.Status.MarkRollbackFailed(rolloutsv1alpha1.RollbackTimedOutReason, "Rollback failed: %v", err)
			return nil
		}
		return err
	}

	if err := r.progressRollout(ctx, rollout, productionDeployment); err != nil {
		logger.Debug("Progressing rollout")
		if errors.Is(err, deployments.ErrReadinessDeadlineExceeded) {
			rollout.Status.MarkFailed(rolloutsv1alpha1.RolloutTimedOutReason, "Rollout of revision %s failed: %v",
				rollout.Spec.Revision, err)
			return r.rollBack(ctx, rollout, productionDeployment)
		}
		return err
	}

	return nil
}

func (r *Reconciler) getProductionDeployment(rollout *rolloutsv1alpha1.Rollout) (deployment *appsv1.Deployment, err error) {
	if rollout.Status.PreviousRollout != nil {
		deployment, err = r.deploymentLister.Deployments(rollout.Namespace).Get(rollout.Status.PreviousRollout.DeploymentName)
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
	} else {
		deployments, err := r.deploymentLister.Deployments(rollout.Namespace).List(labels.SelectorFromSet(labels.Set{
			rolloutsv1alpha1.RolloutAppNameLabelKey: rollout.Spec.AppName,
			rolloutsv1alpha1.RolloutScopeLabelKey:   deployments.ScopeProduction,
		}))

		if err != nil {
			return nil, fmt.Errorf("error listing Deployments from the indexer: %w", err)
		}

		for _, candidate := range deployments {
			if candidate.Name != rollout.Name {
				deployment = candidate
				break
			}
		}
	}
	return
}

func (r *Reconciler) progressRollout(ctx context.Context, rollout *rolloutsv1alpha1.Rollout, productionDeployment *appsv1.Deployment) error {
	logger := logging.FromContext(ctx)

	if rollout.Status.StartedAt == nil {
		rollout.Status.StartedAt = &metav1.Time{Time: clock.Now()}
	}

	if productionDeployment == nil {
		return r.createFirstDeployment(ctx, rollout)
	}

	if rollout.Spec.Canary != nil && !rollout.Status.Canary.Iscompleted() {
		if done, err := r.canary.ProgressRollout(ctx, rollout, productionDeployment); err != nil {
			return err
		} else if !done {
			logger.Debug("Canary's setup hasn't completed yet - waiting for the next reconciliation cycle")
			return nil
		}

		analysis, err := r.ensureAnalysis(ctx, rollout)
		if err != nil {
			return err
		}

		condition := analysis.Status.GetCondition(rolloutsv1alpha1.AnalysisConditionSucceeded)
		if condition.IsUnknown() {
			message := condition.GetMessage()
			if message == "" {
				message = "waiting for analysis to start running"
			}
			rollout.Status.MarkUnknown(rolloutsv1alpha1.RolloutInProgressReason, "Rolling out revision %s (analyzing canary): %s", rollout.Spec.Revision, message)
			return nil
		}

		rollout.Status.Canary.Completed = true

		if condition.IsFalse() {
			rollout.Status.CompletedAt = &metav1.Time{Time: clock.Now()}
			rollout.Status.MarkFailed(rolloutsv1alpha1.RolloutFailedReason, "Rollout of revision %s failed: %s", rollout.Spec.Revision, condition.GetMessage())
			rollout.Status.MarkRollbackUnknown(rolloutsv1alpha1.RollbackInProgressReason, "Rolling back the canary deployment")
			return r.rollBack(ctx, rollout, productionDeployment)
		}
	}

	return r.completeRollout(ctx, rollout, productionDeployment)
}

func (r *Reconciler) createFirstDeployment(ctx context.Context, rollout *rolloutsv1alpha1.Rollout) error {
	desiredReplicas := *rollout.Spec.Replicas

	deployment, err := r.deploymentLister.Deployments(rollout.Namespace).Get(rollout.Name)
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("error reading Deployment %s from the indexer: %w", rollout.Name, err)
	}

	if deployment == nil {
		deployment, err = deployments.Create(ctx, deployments.NewProduction(rollout, desiredReplicas))
	}

	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			err = controller.NewRequeueImmediately()
		}
		return err
	}

	if ready, err := deployments.IsReady(deployment, desiredReplicas); err != nil {
		return err
	} else if !ready {
		rollout.Status.MarkUnknown(rolloutsv1alpha1.RolloutInProgressReason, "Rolling out revision %s (creating Deployment %s): waiting for the underlying Deployment to meet the desired conditions (%d of %d replicas are available)",
			rollout.Spec.Revision, deployment.Name, deployment.Status.AvailableReplicas, desiredReplicas)
		return controller.NewRequeueAfter(defaultRequeueDuration)
	}

	rollout.Status.CompletedAt = &metav1.Time{Time: clock.Now()}
	rollout.Status.MarkSucceeded(rolloutsv1alpha1.RolloutSucceededReason, "Rollout of revision %s completed successfully - Deployment %s is ready", rollout.Spec.Revision, deployment.Name)

	return nil
}

func (r *Reconciler) ensureAnalysis(ctx context.Context, rollout *rolloutsv1alpha1.Rollout) (analysis *rolloutsv1alpha1.Analysis, err error) {
	analysis, err = r.analysisLister.Analysises(rollout.Namespace).Get(rollout.Name)
	if err != nil && !apierrors.IsNotFound(err) {
		return
	}

	if analysis == nil {
		analysis = &rolloutsv1alpha1.Analysis{
			ObjectMeta: metav1.ObjectMeta{
				Name:      rollout.Name,
				Namespace: rollout.Namespace,
				OwnerReferences: []metav1.OwnerReference{
					*kmeta.NewControllerRef(rollout),
				},
			},
			Spec: rolloutsv1alpha1.AnalysisSpec{
				AppName:    rollout.Spec.AppName,
				Duration:   rollout.Spec.Canary.Duration,
				Interval:   rollout.Spec.Canary.Interval,
				RolloutRef: rollout.Name,
				Canary: &rolloutsv1alpha1.CanaryAnalysis{
					MetricGroups: []rolloutsv1alpha1.MetricGroup{{
						Name:   "http-health",
						Weight: pointer.Float64(100.0),
						Metrics: []rolloutsv1alpha1.Metric{{
							Name:     "http-errors",
							Query:    `sum(increase(http_response_total{code!~"2..",pod=~"^${scope}.*"}[20m]))`,
							FailOn:   rolloutsv1alpha1.DirectionIncrease,
							Critical: pointer.Bool(true),
						},
						},
					},
					},
				},
			},
		}

		analysis, err = r.hringhorniClient.HringhorniV1alpha1().Analysises(analysis.Namespace).Create(ctx, analysis, metav1.CreateOptions{})
		if apierrors.IsAlreadyExists(err) {
			rollout.Status.MarkUnknown(rolloutsv1alpha1.RolloutInProgressReason, "Waiting for Analysis %s to be available in the indexer", rollout.Name)
			err = controller.NewRequeueImmediately()
		}
	}
	return
}

func (r *Reconciler) completeRollout(ctx context.Context, rollout *rolloutsv1alpha1.Rollout, productionDeployment *appsv1.Deployment) error {
	logger := logging.FromContext(ctx)

	// First, attempt to revert the canary setup if applicable.
	if done, err := r.undoCanarySetup(ctx, rollout, productionDeployment); err != nil {
		return err
	} else if !done {
		// Wait for the next reconciliation cycle.
		return nil
	}

	if rolloutName := productionDeployment.Labels[rolloutsv1alpha1.RolloutNameLabelKey]; rollout.Name != rolloutName {
		logger.Infof("Deployment %s is out of date - updating", productionDeployment.Name)
		return r.updateDeployment(ctx, rollout, productionDeployment)
	}

	if ready, err := deployments.IsReady(productionDeployment, *rollout.Spec.Replicas); err != nil {
		// Signal to the rollBack method that the object need to be
		// rolled back to the previous version. Note that this will not
		// be reflected in the Rollout's actual Spec field in the
		// cluster.
		rollout.Spec.Actions.RollBack = true
		return err
	} else if !ready {
		rollout.Status.MarkUnknown(rolloutsv1alpha1.RolloutInProgressReason, "Rolling out revision %s (completing rollout): waiting for Deployment %s to be ready (%d of %d desired replicas are available)",
			rollout.Spec.Revision, productionDeployment.Name, productionDeployment.Status.AvailableReplicas, *rollout.Spec.Replicas)
		return controller.NewRequeueAfter(defaultRequeueDuration)
	}

	rollout.Status.CompletedAt = &metav1.Time{Time: clock.Now()}
	rollout.Status.MarkSucceeded(rolloutsv1alpha1.RolloutSucceededReason, "Rollout of revision %s completed successfully", rollout.Spec.Revision)
	return nil
}

func (r *Reconciler) updateDeployment(ctx context.Context, rollout *rolloutsv1alpha1.Rollout, productionDeployment *appsv1.Deployment) error {
	newDeployment := deployments.NewProduction(rollout, *rollout.Spec.Replicas)
	newDeployment.Name = productionDeployment.Name
	newDeployment.ResourceVersion = productionDeployment.ResourceVersion
	if err := deployments.Update(ctx, newDeployment); err != nil {
		return err
	}
	rollout.Status.MarkUnknown(rolloutsv1alpha1.RolloutInProgressReason, "Rolling out revision %s (completing rollout): updating Deployment %s to match the desired configurations",
		rollout.Spec.Revision, newDeployment.Name)
	return controller.NewRequeueAfter(defaultRequeueDuration)
}
