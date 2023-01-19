package test

/*
Copyright 2022 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	"testing"
	"time"

	"github.com/k8gb-io/k8gb-light/terratest"
	"github.com/k8gb-io/k8gb-light/terratest/utils"

	"github.com/stretchr/testify/assert"
)

func TestDNSEndpointLifecycle(t *testing.T) {
	const ingressPath = "./resources/ingress_rr_ep_lifecycle.yaml"
	const ingressEmptyPath = "./resources/ingress_empty_ep_lifecycle.yaml"
	const ingressInvalidPath = "./resources/ingress_invalid_ep_lifecycle.yaml"
	instanceEU, err := utils.NewWorkflow(t, terratest.Environment.EUCluster, terratest.Environment.EUClusterPort).
		WithIngress(ingressPath).
		WithTestApp("eu").
		Start()
	assert.NoError(t, err)
	defer instanceEU.Kill()
	info := instanceEU.GetInfo()

	t.Run("Apply ingress with k8gb annotation, DNSEndpoint created", func(t *testing.T) {
		err = instanceEU.Resources().WaitUntilDNSEndpointContainsTargets(info.Host, info.NodeIPs)
		assert.NoError(t, err)
	})

	t.Run("Remove k8gb annotation from ingress, DNSEndpoint removed", func(t *testing.T) {
		instanceEU.ReapplyIngress(ingressEmptyPath)
		err = instanceEU.Resources().WaitUntilDNSEndpointNotFound()
		assert.NoError(t, err)
	})

	t.Run("Applying ingress with invalid strategy; DNSEndpoint removed", func(t *testing.T) {
		instanceEU.ReapplyIngress(ingressInvalidPath)
		time.Sleep(5 * time.Second)
		err = instanceEU.Resources().WaitUntilDNSEndpointNotFound()
		assert.NoError(t, err)
	})
}
