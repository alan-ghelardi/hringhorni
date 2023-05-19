package rollout

import (
	"context"
	"time"

	hringhorniclient "github.com/nubank/hringhorni/pkg/client/injection/client"
	"github.com/nubank/hringhorni/pkg/rolloutstrategy"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"

	analysisinformer "github.com/nubank/hringhorni/pkg/client/injection/informers/rollouts/v1alpha1/analysis"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	deploymentinformer "knative.dev/pkg/client/injection/kube/informers/apps/v1/deployment"
	"knative.dev/pkg/kmeta"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"

	"github.com/nubank/hringhorni/pkg/apis/config"
	rolloutsv1alpha1 "github.com/nubank/hringhorni/pkg/apis/rollouts/v1alpha1"
	rolloutinformer "github.com/nubank/hringhorni/pkg/client/injection/informers/rollouts/v1alpha1/rollout"
	rolloutreconciler "github.com/nubank/hringhorni/pkg/client/injection/reconciler/rollouts/v1alpha1/rollout"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/logging"
)

var defaultDelay = 30 * time.Second

// NewController creates a Reconciler and returns the result of NewImpl.
func NewController(ctx context.Context, watcher configmap.Watcher) *controller.Impl {
	logger := logging.FromContext(ctx)

	configStore := config.NewStore(logger.Named("configs"))
	configStore.WatchConfigs(watcher)

	rolloutInformer := rolloutinformer.Get(ctx)
	analysisInformer := analysisinformer.Get(ctx)
	deploymentInformer := deploymentinformer.Get(ctx)

	reconciler := &Reconciler{
		analysisLister:   analysisInformer.Lister(),
		canary:           rolloutstrategy.New(ctx),
		deploymentLister: deploymentInformer.Lister(),
		hringhorniClient: hringhorniclient.Get(ctx),
		kubeClient:       kubeclient.Get(ctx),
	}

	impl := rolloutreconciler.NewImpl(ctx, reconciler, func(*controller.Impl) controller.Options {
		return controller.Options{ConfigStore: configStore}
	})

	logger.Info("Setting up event handlers")

	rolloutInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	eventHandler := cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterController(&rolloutsv1alpha1.Rollout{}),
		Handler:    controller.HandleAll(enqueueControllerOfAfterDelay(ctx, impl, defaultDelay)),
	}

	analysisInformer.Informer().AddEventHandler(eventHandler)

	deploymentInformer.Informer().AddEventHandler(eventHandler)

	return impl
}

func enqueueControllerOfAfterDelay(ctx context.Context, impl *controller.Impl, delay time.Duration) controller.Callback {
	return func(in any) {
		object, err := kmeta.DeletionHandlingAccessor(in)
		if err != nil {
			logger := logging.FromContext(ctx)
			logger.Errorw("Error processing object", zap.Error(err))
			return
		}

		if owner := metav1.GetControllerOf(object); owner != nil {
			impl.EnqueueKeyAfter(types.NamespacedName{
				Namespace: object.GetNamespace(),
				Name:      owner.Name,
			}, delay)
		}
	}
}
