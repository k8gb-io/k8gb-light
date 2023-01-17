TERRATEST_DIR ?=$(CURDIR)

DEFAULT_CONTEXT ?=k3d-k8gb-test-eu

deploy-demo:
	$(MAKEIN) deploy-demo-k8gb-test-eu
	$(MAKEIN) deploy-demo-k8gb-test-us
	$(MAKEIN) deploy-demo-k8gb-test-za
	kubectl config use-context $(DEFAULT_CONTEXT)
	kubectl config set-context --current --namespace=demo
	clear

deploy-demo-%:
	kubectl delete namespace demo --context=k3d-$(*) --ignore-not-found=true
	kubectl create namespace demo --dry-run=client -o yaml | kubectl apply --force -f - --context=k3d-$(*)
	helm upgrade --install frontend podinfo/podinfo \
		--set ui.message="k3d-$(*)" \
		--set ui.color="#34577c"\
		--set image.repository=ghcr.io/stefanprodan/podinfo \
		--version 5.2.0 \
		--kube-context k3d-$(*)  \
		--namespace demo

init-failover:
	kubectl apply -f $(TERRATEST_DIR)/deploy/demo/fo_demo_ingress.yaml -n demo --context=k3d-k8gb-test-eu
	kubectl apply -f $(TERRATEST_DIR)/deploy/demo/fo_demo_ingress.yaml -n demo --context=k3d-k8gb-test-us
	kubectl apply -f $(TERRATEST_DIR)/deploy/demo/fo_demo_ingress.yaml -n demo --context=k3d-k8gb-test-za

init-wrr:
	kubectl apply -f $(TERRATEST_DIR)/deploy/demo/wrr_demo_ingress.yaml -n demo --context=k3d-k8gb-test-eu
	kubectl apply -f $(TERRATEST_DIR)/deploy/demo/wrr_demo_ingress.yaml -n demo --context=k3d-k8gb-test-us
	kubectl apply -f $(TERRATEST_DIR)/deploy/demo/wrr_demo_ingress.yaml -n demo --context=k3d-k8gb-test-za

kill-local-k8gb:
	kubectl config use-context $(DEFAULT_CONTEXT)
	kubectl config set-context --current --namespace=demo
	kubectl -n k8gb scale deployment k8gb --replicas=0 --context=$(DEFAULT_CONTEXT)

start-local-k8gb:
	kubectl config use-context $(DEFAULT_CONTEXT)
	kubectl config set-context --current --namespace=demo
	kubectl -n k8gb scale deployment k8gb --replicas=1 --context=$(DEFAULT_CONTEXT)

I ?=1
dig:
	@for run in {1..$(I)}; do dig -p 5063 @localhost demo.cloud.example.com -4 +tcp +short; done