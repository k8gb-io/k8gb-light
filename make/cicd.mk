# cicd.mk contains
BIN := k8gb
REGISTRY = docker.io
REPOSITORY = kuritka
TAG = 932-fix-3
IMG = $(REGISTRY)/$(REPOSITORY)/$(BIN)
OPTIONAL_DOCKERFILE_PATH ?= .

image: build
	docker build $(OPTIONAL_DOCKERFILE_PATH) -t $(IMG):$(TAG)

rebuild-k8gb-image:
	@echo -e "\n$(YELLOW)Rebuild Image $(CYAN)$(REPOSITORY)/$(BIN):$(TAG) $(NC)"
	$(MAKEIN) image

push:
	docker push $(IMG):$(TAG)