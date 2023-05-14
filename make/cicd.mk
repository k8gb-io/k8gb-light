# cicd.mk contains
BIN := k8gb
REGISTRY = docker.io
REPOSITORY = kuritka
TAG = 932-fix-3
IMG = $(REGISTRY)/$(REPOSITORY)/$(BIN)

image: build
	docker build . -t ${IMG}:${TAG}
