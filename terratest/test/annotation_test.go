package test

import (
	"github.com/kuritka/annotation-operator/terratest/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAnnotation(t *testing.T) {
	const ingressPath1 = "./resources/ingress_annotation1.yaml"
	const ingressPath2 = "./resources/ingress_annotation2.yaml"
	const ingressPath3 = "./resources/ingress_annotation3.yaml"
	const ingressPath4 = "./resources/ingress_annotation4.yaml"
	instanceEU, err := utils.NewWorkflow(t, ContextEU, PortEU).
		WithIngress(ingressPath1).
		WithTestApp("eu").
		Start()
	assert.NoError(t, err)
	defer instanceEU.Kill()

	t.Run("Wait until cluster is ready", func(t *testing.T) {
		euClusterIPs := instanceEU.GetInfo().NodeIPs
		err = instanceEU.Resources().WaitUntilDNSEndpointContainsTargets(instanceEU.GetInfo().Host, euClusterIPs)
		require.NoError(t, err)
	})
}
