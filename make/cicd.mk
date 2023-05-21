# cicd.mk contains
BIN := k8gb
REGISTRY = docker.io
REPOSITORY = kuritka
TAG = 932-fix-3
IMG = $(REGISTRY)/$(REPOSITORY)/$(BIN)
OPTIONAL_DOCKERFILE_PATH ?= .

image: build
	docker build . $(OPTIONAL_DOCKERFILE_PATH) -t ${IMG}:${TAG}
