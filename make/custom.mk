
ZSCALER_CERT_PATH ?= $(shell pwd ./zscaler.crt)/zscaler.crt
OPTIONAL_CLUSTER_ARGUMENTS=--volume $(ZSCALER_CERT_PATH):/etc/ssl/certs/zscaler.crt
OPTIONAL_DOCKERFILE_PATH= -f ./custom.Dockerfile
COLIMA_PROFILE ?= k8gb


with-colima:
	colima delete $(COLIMA_PROFILE) && colima start -c 4 -m 8 -p $(COLIMA_PROFILE)
	@echo -e "\n$(YELLOW)Colima certificates$(NC)"
	colima ssh -p $(COLIMA_PROFILE) -- sudo sh -c \
	"openssl s_client -showcerts -connect pkg-containers.githubusercontent.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/pkg-containers.githubusercontent.com.crt && \
	openssl s_client -showcerts -connect ghcr.io:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/ghcr.io.crt && \
	openssl s_client -showcerts -connect k8s.gcr.io:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/k8s.gcr.io.crt && \
	openssl s_client -showcerts -connect gcr.io:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/gcr.io.crt && \
	openssl s_client -showcerts -connect github.io:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/github.io.crt && \
	openssl s_client -showcerts -connect k8gb.io:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/k8gb.io.crt && \
	openssl s_client -showcerts -connect docker.io:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/docker.io.crt && \
	openssl s_client -showcerts -connect production.cloudflare.docker.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/production.cloudflare.docker.com.crt && \
	openssl s_client -showcerts -connect docker.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/docker.com.crt && \
	openssl s_client -showcerts -connect hub.docker.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/hub.docker.com.crt && \
	openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/proxy.golang.org.crt && \
	openssl s_client -showcerts -connect golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/golang.org.crt && \
	openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/github.com.crt && \
	openssl s_client -showcerts -connect golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/golang.org.crt && \
	openssl s_client -showcerts -connect registry.k8s.io:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/registry.k8s.io.crt && \
	openssl s_client -showcerts -connect europe-west4-docker.pkg.dev:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/europe-west4-docker.pkg.dev.crt && \
	openssl s_client -showcerts -connect storage.googleapis.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/storage.googleapis.com.crt && \
	openssl s_client -showcerts -connect production.cloudflare.docker.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/production.cloudflare.docker.com.crt && \
	openssl s_client -showcerts -connect quay.io:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > /usr/local/share/ca-certificates/quay.io.crt && \
	update-ca-certificates"
	colima ssh -p $(COLIMA_PROFILE) -- sudo sh -c "cat /var/run/docker.pid | xargs kill"
	sleep 8
