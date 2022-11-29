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
