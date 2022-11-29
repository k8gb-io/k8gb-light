include Makefile.Gen

ing: deploy-ing
install-ing: deploy-ing
deploy-ing:
	kubectl -n demo apply -f ing.yaml --context=k3d-test-gslb1
	kubectl -n demo apply -f ing.yaml --context=k3d-test-gslb2


MAKEIN=make -C .
KEY?=gslb
watch:
	@watch -n 1 $(MAKEIN) generation

generation:
	@echo "DNSEndpoint generation"  `kubectl get dnsendpoint $(KEY) -ojsonpath={.metadata.generation}`
	@echo "Ingress generation" `kubectl get ingress $(KEY) -ojsonpath={.metadata.generation}`
	@echo
	@kubectl get dnsendpoint $(KEY) -oyaml | grep "  endpoints:" -A 23
