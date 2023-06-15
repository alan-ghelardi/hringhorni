lint:
	@golangci-lint run

install-dev:
	@export KO_DOCKER_REPO=kind.local && kustomize build --enable-helm config/base | ko apply -f -
