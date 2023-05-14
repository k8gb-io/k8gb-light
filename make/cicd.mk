zBIN := k8gb
REGISTRY = docker.io
REPOSITORY = kuritka
TAG = 932-fix-3
IMG = $(REGISTRY)/$(REPOSITORY)/$(BIN)

image:
	docker build . -t ${IMG}:${TAG}
