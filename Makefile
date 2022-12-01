BIN := k8gb
REGISTRY = docker.io
REPOSITORY = kuritka
TAG = 932-fix-3
IMG = $(REGISTRY)/$(REPOSITORY)/$(BIN)
KEY?=ing
NS?=demo

watch:
	@watch -n 1 $(MAKEIN) generation

generation:
	@echo "DNSEndpoint generation"  `kubectl get dnsendpoint $(KEY) -ojsonpath={.metadata.generation} -n $(NS)`
	@echo "Ingress generation" `kubectl get ingress $(KEY) -ojsonpath={.metadata.generation} -n $(NS)`
	@echo
	@kubectl get dnsendpoint $(KEY) -oyaml -n $(NS) | grep "  endpoints:" -A 23


ing:
	kubectl -n demo apply -f ing.yaml --context=k3d-test-gslb1
	kubectl -n demo apply -f ing.yaml --context=k3d-test-gslb2

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./k8gb main.go

image:
	docker build . -t ${IMG}:${TAG}


redeploy: build image
	k3d image import ${REPOSITORY}/${BIN}:${TAG} -c test-gslb2
	k3d image import ${REPOSITORY}/${BIN}:${TAG} -c test-gslb1
	kubectl -n k8gb patch deployment k8gb -p '{"spec": {"template":{"spec":{"containers":[{"name":"k8gb","image":"$(IMG):$(TAG)"}]}}}}'
