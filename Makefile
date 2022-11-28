include Makefile.Gen

ing: deploy-ing
install-ing: deploy-ing
deploy-ing:
	kubectl -n demo apply -f ing.yaml --context=k3d-test-gslb1
	kubectl -n demo apply -f ing.yaml --context=k3d-test-gslb2

