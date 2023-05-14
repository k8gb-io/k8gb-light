MAKEIN ?=make -C .

NGINX_INGRESS_VALUES_PATH ?= $(TERRATEST_DIR)/deploy/ingress/nginx-ingress-values.yaml
CLUSTER_GSLB_NETWORK = k3d-action-bridge-network
CLUSTER_GSLB_GATEWAY = docker network inspect $(CLUSTER_GSLB_NETWORK) -f '{{ (index .IPAM.Config 0).Gateway }}'

create-clusters:
	$(MAKEIN) prepare-cluster-edge-dns
	$(MAKEIN) prepare-cluster-k8gb-test-eu
	$(MAKEIN) prepare-cluster-k8gb-test-us
	$(MAKEIN) prepare-cluster-k8gb-test-za