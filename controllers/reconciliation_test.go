package controllers

import (
	"cloud.example.com/annotation-operator/controllers/reconciliation"
	"cloud.example.com/annotation-operator/controllers/utils"
	"k8s.io/component-base/metrics/prometheus/controllers"
	"sigs.k8s.io/controller-runtime/pkg/config"
	"testing"
	//"github.com/stretchr/testify/require"
	//"github.com/stretchr/testify/assert"
)

func TestReconcile(t *testing.T) {

	reconciler := &controllers.AnnoReconciler{
		Config:           config,
		Client:           mgr.GetClient(),
		DepResolver:      resolver,
		Scheme:           mgr.GetScheme(),
		IngressMapper:    reconciliation.NewIngressMapper(mgr.GetClient()),
		ReconcilerResult: utils.NewReconcileResultHandler(config.ReconcileRequeueSeconds),
	}
}
