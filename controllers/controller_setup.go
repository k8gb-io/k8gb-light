package controllers

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

// SetupWithManager sets up the controller with the Manager.
func (r *AnnoReconciler) SetupWithManager(mgr ctrl.Manager) error {

	ingressHandler := handler.EnqueueRequestsFromMapFunc(
		func(a client.Object) []reconcile.Request {
			// thanks to watcher, the reconciliation is executed immediatelly at the moment when ingress changed
			// skip
			return nil
		})

	serviceEndpointHandler := handler.EnqueueRequestsFromMapFunc(
		func(a client.Object) []reconcile.Request {
			ingList := &netv1.IngressList{}
			c := mgr.GetClient()
			err := c.List(context.TODO(), ingList, client.InNamespace(a.GetNamespace()))
			if err != nil {
				log.Info().Msg("Can't fetch ingress objects")
				return nil
			}
			for _, ing := range ingList.Items {
				for _, rule := range ing.Spec.Rules {
					for _, path := range rule.HTTP.Paths {
						if path.Backend.Service != nil && path.Backend.Service.Name == a.GetName() {
							return []reconcile.Request{{types.NamespacedName{Namespace: a.GetNamespace(), Name: ing.Name}}}
						}
					}
				}
			}
			return nil
		})

	return ctrl.NewControllerManagedBy(mgr).
		For(&netv1.Ingress{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Owns(&externaldns.DNSEndpoint{}).
		Watches(&source.Kind{Type: &netv1.Ingress{}}, ingressHandler).
		Watches(&source.Kind{Type: &corev1.Endpoints{}}, serviceEndpointHandler).
		Complete(r)
}
