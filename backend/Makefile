DOCKER_IMAGE := local/mex

GO_VERSION := 1.20
CILINT_VERSION := v1.52

unit-test:
	docker build \
		--build-arg GO_VERSION="$(GO_VERSION)" \
		--target unit-test \
		-t "$(DOCKER_IMAGE):test" \
		.

lint:
	DOCKER_BUILDKIT=1 \
	docker build \
		--build-arg GO_VERSION="$(GO_VERSION)" \
		--build-arg CILINT_VERSION=${CILINT_VERSION} \
		--target lint \
		-t "$(DOCKER_IMAGE):lint \
		.
