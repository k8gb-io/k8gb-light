KEY?=ing
NS?=demo

CTX?=k3d-k8gb-test-eu
watch:
	@watch -n 1 $(MAKEIN) generation

generation:
	@echo "CONTEXT $(CTX):" `kubectl get deployment -n k8gb k8gb -ojsonpath='{.spec.template.spec.containers[*].env[?(@.name=="CLUSTER_GEO_TAG")].value}' --context=$(CTX)`
	@echo "DNSEndpoint generation:"  `kubectl get dnsendpoint $(KEY) -ojsonpath={.metadata.generation} -n $(NS) --context=$(CTX) 2>/dev/null`
	@echo "Ingress generation:" `kubectl get ingress $(KEY) -ojsonpath={.metadata.generation} -n $(NS) --context=$(CTX) 2>/dev/null`
	@echo "Finalizers:" `kubectl get ingress $(KEY) -ojsonpath={.metadata.finalizers} -n $(NS) --context=$(CTX) 2>/dev/null`
	@echo "Image" `kubectl get deployment  -n k8gb k8gb -oyaml  --context=$(CTX) | grep "        image:"`
	@echo "primary-geo-tag:" `kubectl get ing ing -ojsonpath='{.metadata.annotations.k8gb\.io/primary-geotag}' --context=$(CTX) 2>/dev/null`
	@echo "status:" `kubectl get ing ing -ojsonpath='{.metadata.annotations.k8gb\.io/status}' -n $(NS) --context=$(CTX) 2>/dev/null`
	@echo
	@kubectl get dnsendpoint $(KEY) -oyaml -n $(NS) 2>/dev/null --context=$(CTX) | grep "  endpoints:" -A 23
