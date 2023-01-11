# Copyright 2022 The k8gb Contributors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
BIN := k8gb
REGISTRY = docker.io
REPOSITORY = kuritka
TAG = 932-fix-3
IMG = $(REGISTRY)/$(REPOSITORY)/$(BIN)
KEY?=ing
NS?=demo

GOLIC_VERSION  ?= v0.7.2
GOKART_VERSION ?= v0.5.1
GOLANGCI_VERSION ?= v1.50.1
MOCKGEN_VERSION ?= v1.6.0

SHELL := bash
TERRATEST_DIR =$(CURDIR)/terratest
MAKEIN =make -C .
MAKEAWAY =make -C .

ifndef NO_COLOR
YELLOW=\033[0;33m
CYAN=\033[1;36m
RED=\033[31m
# no color
NC=\033[0m
endif

# create GOBIN if not specified
ifndef GOBIN
GOBIN=$(shell go env GOPATH)/bin
endif

# check integrity
.PHONY: quick-check
quick-check: lint test ## Check project integrity

.PHONY: check
check: mocks gokart build quick-check  ## Check project integrity

# updates source code with license headers
.PHONY: license
license:
	@echo -e "\n$(YELLOW)Injecting the license$(NC)"
	@go install github.com/AbsaOSS/golic@$(GOLIC_VERSION)
	$(GOBIN)/golic inject -t apache2

# runs golangci-lint aggregated linter; see .golangci.yaml for linter list
.PHONY: lint
lint:
	@echo -e "\n$(YELLOW)Running the linters$(NC)"
	goimports -w ./
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_VERSION)
	$(GOBIN)/golangci-lint run

.PHONY: build
build:
	@echo -e "\n$(YELLOW)Building binary$(NC)"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./k8gb main.go

# run tests
.PHONY: test
test:
	@echo -e "\n$(YELLOW)Running the unit tests$(NC)"
	go test $$(go list ./... | grep -Ev '/mocks|/terratest|/logging|/tracing') --cover

# GoKart - Go Security Static Analysis
# see: https://github.com/praetorian-inc/gokart
.PHONY: gokart
gokart:
	@go install github.com/praetorian-inc/gokart@$(GOKART_VERSION)
	$(GOBIN)/gokart scan --globalsTainted --verbose

.PHONY: mocks
mocks:
	go install github.com/golang/mock/mockgen@$(MOCKGEN_VERSION)
	mockgen -package=mocks -destination=controllers/mocks/assistant_mock.go -source=controllers/providers/assistant/assistant.go Assistant
	mockgen -package=mocks -destination=controllers/mocks/client_mock.go sigs.k8s.io/controller-runtime/pkg/client Client
	mockgen -package=mocks -destination=controllers/mocks/resolver_mock.go -source=controllers/depresolver/resolver.go GslbResolver
	mockgen -package=mocks -destination=controllers/mocks/provider_mock.go -source=controllers/providers/dns/dns.go Provider
	mockgen -package=mocks -destination=controllers/mocks/mapper_mock.go -source=controllers/mapper/mapper.go Mapper
	mockgen -package=mocks -destination=controllers/mocks/mapper_provider_mock.go -source=controllers/mapper/provider.go ProviderMapper
	mockgen -package=mocks -destination=controllers/mocks/manager_mock.go sigs.k8s.io/controller-runtime/pkg/manager Manager
	mockgen -package=mocks -destination=controllers/mocks/infoblox-client_mock.go -source=controllers/providers/dns/infoblox-client.go InfobloxClient
	mockgen -package=mocks -destination=controllers/mocks/infoblox-connection_mock.go github.com/infobloxopen/infoblox-go-client IBConnector
	mockgen -package=mocks -destination=controllers/mocks/tracer_mock.go go.opentelemetry.io/otel/trace Tracer
	mockgen -package=mocks -destination=controllers/mocks/span_mock.go go.opentelemetry.io/otel/trace Span
	mockgen -package=mocks -destination=controllers/mocks/metrics_mock.go -source=controllers/providers/metrics/provider.go Provider
	mockgen -package=mapper -destination=controllers/mapper/dig_mock.go -source=controllers/utils/dns.go Digger
	mockgen -package=mapper -destination=controllers/mapper/client_mock.go sigs.k8s.io/controller-runtime/pkg/client Client
	$(MAKEIN) license

image:
	docker build . -t ${IMG}:${TAG}

include ./terratest/Makefile