# local-platform.mk is all about provisioning local clusters

# OPTIONAL_CLUSTER_ARGUMENTS exists to be overridden in commandline or in custom.mk (custom.mk is .gitignored).
# The purpose is to be able to inject custom arguments to the clusters. See example injecting certificate:
# ./make/custom.mk
# CERT_PATH ?= $(shell pwd ./cert.crt)/cert.crt
# OPTIONAL_CLUSTER_ARGUMENTS = --volume $(CERT_PATH):/etc/ssl/certs/cert.crt
OPTIONAL_CLUSTER_ARGUMENTS ?=

NGINX_INGRESS_VALUES_PATH ?= $(TERRATEST_DIR)/deploy/ingress/nginx-ingress-values.yaml
CLUSTER_GSLB_NETWORK = k3d-action-bridge-network
CLUSTER_GSLB_GATEWAY = docker network inspect $(CLUSTER_GSLB_NETWORK) -f '{{ (index .IPAM.Config 0).Gateway }}'


deploy-full-local-setup:
	$(MAKEIN) create-clusters
	$(MAKEIN) deploy-full-terratest-setup
	$(MAKEIN) deploy-demo


deploy-full-terratest-setup: image
	$(MAKEIN) deploy-cluster-k8gb-edge-dns
	$(MAKEIN) deploy-cluster-k8gb-test-eu CLUSTER_GEOTAG=eu EXT_CLUSTER_GEOTAGS=us\\,za
	$(MAKEIN) deploy-cluster-k8gb-test-us CLUSTER_GEOTAG=us EXT_CLUSTER_GEOTAGS=eu\\,za
	$(MAKEIN) deploy-cluster-k8gb-test-za CLUSTER_GEOTAG=za EXT_CLUSTER_GEOTAGS=eu\\,us

create-clusters:
	$(MAKEIN) prepare-cluster-edge-dns
	$(MAKEIN) prepare-cluster-k8gb-test-eu
	$(MAKEIN) prepare-cluster-k8gb-test-us
	$(MAKEIN) prepare-cluster-k8gb-test-za

prepare-cluster-%:
	@echo -e "\n$(YELLOW)Prepare cluster $(CYAN)$(*)$(NC)"
	k3d cluster delete $(*)
	k3d cluster create -c $(TERRATEST_DIR)/deploy/k3d/$(*).yaml $(*) $(OPTIONAL_CLUSTER_ARGUMENTS)

deploy-cluster-k8gb-edge-dns:
	@echo -e "\n$(YELLOW)Deploying EdgeDNS $(NC)"
	kubectl --context k3d-edge-dns apply -f $(TERRATEST_DIR)/deploy/edge/

deploy-cluster-%:
	@echo -e "\n$(YELLOW)Deploy k8gb on cluster $(CYAN)$(*)$(NC)"
	@echo -e "\n$(YELLOW)Create namespace $(NC)"
	kubectl delete namespace k8gb --context=k3d-$(*) --ignore-not-found
	kubectl create namespace k8gb --context=k3d-$(*)

	@echo -e "\n$(YELLOW)Deploy Ingress $(NC)"
	helm repo add --force-update nginx-stable https://kubernetes.github.io/ingress-nginx
	helm repo add --force-update k8gb https://www.k8gb.io
	helm repo update
	helm -n k8gb upgrade -i nginx-ingress nginx-stable/ingress-nginx \
		--version 4.0.15 -f $(NGINX_INGRESS_VALUES_PATH) --kube-context=k3d-$(*)

	@echo -e "\n$(YELLOW)Deploy K8GB $(NC)"
	kubectl -n k8gb --context=k3d-$(*) create secret generic rfc2136 --from-literal=secret=96Ah/a2g0/nLeFGK+d/0tzQcccf9hCEIy34PoXX2Qg8= || true
	cd $(TERRATEST_DIR)/deploy/chart/k8gb && helm dependency update
	helm -n k8gb upgrade -i k8gb k8gb/k8gb \
		--set k8gb.clusterGeoTag='$(CLUSTER_GEOTAG)' \
		--set k8gb.extGslbClustersGeoTags='$(EXT_CLUSTER_GEOTAGS)' \
		--set k8gb.reconcileRequeueSeconds=10 \
		--set k8gb.dnsZoneNegTTL=10 \
		--set k8gb.imageTag=${VERSION:"stable"=""} \
		--set k8gb.log.format=simple \
		--set k8gb.log.level=debug \
		--set rfc2136.enabled=true \
		--set k8gb.edgeDNSServers[0]=$(shell $(CLUSTER_GSLB_GATEWAY)):1053 \
		--set externaldns.image=absaoss/external-dns:rfc-ns1 \
		--wait --timeout=2m0s --kube-context=k3d-$(*)

	@echo -e "\n$(YELLOW)Deploy CoreDNS $(NC)"
	kubectl apply -f $(TERRATEST_DIR)/deploy/coredns --context=k3d-$(*)

	@echo -e "\n$(YELLOW)Patch K8GB to local version $(NC)"
	kubectl -n k8gb patch deployment k8gb -p '{"spec": {"template":{"spec":{"containers":[{"name":"k8gb","image":"$(IMG):$(TAG)"}]}}}}' --context=k3d-$(*)
