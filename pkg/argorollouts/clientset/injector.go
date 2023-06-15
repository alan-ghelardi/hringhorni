package clientset

import (
	"context"

	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/logging"

	clientset "github.com/argoproj/argo-rollouts/pkg/client/clientset/versioned"
)

func init() {
	injection.Default.RegisterClient(func(ctx context.Context, config *rest.Config) context.Context {
		return withClient(ctx, config)
	})
}

type clientKey struct {
}

func Get(ctx context.Context) clientset.Interface {
	client, ok := ctx.Value(clientKey{}).(clientset.Interface)
	if !ok {
		logger := logging.FromContext(ctx)
		logger.Panic("Unable to fetch a Argo Rollouts client set from the provided Context")
	}
	return client
}

func withClient(ctx context.Context, config *rest.Config) context.Context {
	client, err := clientset.NewForConfig(config)
	if err != nil {
		logger := logging.FromContext(ctx)
		logger.Panicw("error creating Argo Rollouts client set", zap.Error(err))
	}
	return context.WithValue(ctx, clientKey{}, client)
}
