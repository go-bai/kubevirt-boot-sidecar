export DOCKER_PREFIX=ghcr.io/go-bai
export DOCKER_TAG=v1.2.0

push:
	docker build --tag ${DOCKER_PREFIX}/kubevirt-boot-sidecar:${DOCKER_TAG} --file Dockerfile --push .

.PHONY: \
	push