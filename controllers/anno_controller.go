/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"cloud.example.com/annotation-operator/controllers/depresolver"
	"cloud.example.com/annotation-operator/controllers/providers/metrics"
	"cloud.example.com/annotation-operator/controllers/rs"
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	externaldns "sigs.k8s.io/external-dns/endpoint"

	"cloud.example.com/annotation-operator/controllers/logging"
	"cloud.example.com/annotation-operator/controllers/providers/dns"
	"go.opentelemetry.io/otel/trace"

	"cloud.example.com/annotation-operator/controllers/utils"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// AnnoReconciler reconciles a Anno object
type AnnoReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	Config      *depresolver.Config
	DepResolver depresolver.GslbResolver
	DNSProvider dns.Provider
	Tracer      trace.Tracer
}

var log = logging.Logger()

var m = metrics.Metrics()

func (r *AnnoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log.Info().Msg("Reconciliation started")

	result := utils.NewReconcileResultHandler(30)

	ing := &netv1.Ingress{}
	err := r.Get(ctx, req.NamespacedName, ing)
	if err != nil {
		log.Err(err).Msg("Ingress load error")
	}

	m, err := rs.NewReconciliationState(ing)
	if err != nil {
		log.Err(err).Msg("Invalid ingress")
		return result.Requeue()
	}

	if !m.HasStrategy() {
		log.Info().Msg("No annotation found")
		return result.Requeue()
	}

	return result.Requeue()
}

// SetupWithManager sets up the controller with the Manager.
func (r *AnnoReconciler) SetupWithManager(mgr ctrl.Manager) error {

	ingressHandler := handler.EnqueueRequestsFromMapFunc(
		func(a client.Object) []reconcile.Request {
			// thanks to watcher, the reconciliation is executed immediatelly at the moment when ingress changed
			// skip
			return nil
		})

	endpointHandler := handler.EnqueueRequestsFromMapFunc(
		func(a client.Object) []reconcile.Request {
			var ingress = &netv1.Ingress{}
			c := mgr.GetClient()
			err := c.Get(context.TODO(), client.ObjectKey{
				Namespace: a.GetNamespace(),
				Name:      a.GetName(),
			}, ingress)
			if err == nil {
				return nil
			}
			return []reconcile.Request{{types.NamespacedName{Namespace: ingress.Namespace, Name: ingress.Name}}}
		})

	return ctrl.NewControllerManagedBy(mgr).
		For(&netv1.Ingress{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(&source.Kind{Type: &netv1.Ingress{}}, ingressHandler).
		Owns(&externaldns.DNSEndpoint{}).
		Watches(&source.Kind{Type: &corev1.Endpoints{}}, endpointHandler).
		Complete(r)
}
