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
	"context"

	"cloud.example.com/annotation-operator/controllers/reconciliation"

	"cloud.example.com/annotation-operator/controllers/depresolver"
	"cloud.example.com/annotation-operator/controllers/logging"
	"cloud.example.com/annotation-operator/controllers/providers/dns"
	"cloud.example.com/annotation-operator/controllers/providers/metrics"
	"go.opentelemetry.io/otel/trace"

	"cloud.example.com/annotation-operator/controllers/utils"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AnnoReconciler reconciles a Anno object
type AnnoReconciler struct {
	client.Client
	Scheme           *runtime.Scheme
	Config           *depresolver.Config
	DepResolver      depresolver.GslbResolver
	DNSProvider      dns.Provider
	Tracer           trace.Tracer
	IngressMapper    *reconciliation.IngressMapper
	ReconcilerResult *utils.ReconcileResultHandler
}

var log = logging.Logger()

var m = metrics.Metrics()

func (r *AnnoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx, span := r.Tracer.Start(ctx, "Reconcile")
	defer span.End()

	if req.NamespacedName.Name == "" || req.NamespacedName.Namespace == "" {
		return r.ReconcilerResult.Requeue()
	}
	// TODO: add finalizer for infoblox only
	rs, rr, err := r.IngressMapper.Get(req.NamespacedName)
	switch rr {
	case reconciliation.MapperResultCreate:
		log.Info().
			Str("Namespace", req.NamespacedName.Namespace).
			Str("Ingress", req.NamespacedName.Name).
			Msg("Ingress not found. Stop...")
		return r.ReconcilerResult.Stop()
	case reconciliation.MapperResultError:
		m.IncrementError(rs)
		log.Err(err).
			Str("Namespace", req.NamespacedName.Namespace).
			Str("Ingress", req.NamespacedName.Name).
			Msg("reading Ingress error")
		return r.ReconcilerResult.Requeue()
	}

	if !rs.HasStrategy() {
		log.Info().
			Str("annotation", reconciliation.AnnotationStrategy).
			Msg("No annotation found")
		return r.ReconcilerResult.Requeue()
	}

	log.Info().
		Str("EdgeDNSZone", r.Config.DNSZone).
		Msg("* Starting Reconciliation")

	// == external-dns dnsendpoints CRs ==
	dnsEndpoint, err := r.gslbDNSEndpoint(rs)
	if err != nil {
		m.IncrementError(rs)
		return r.ReconcilerResult.RequeueError(err)
	}

	_, s := r.Tracer.Start(ctx, "SaveDNSEndpoint")
	err = r.DNSProvider.SaveDNSEndpoint(rs, dnsEndpoint)
	if err != nil {
		m.IncrementError(rs)
		return r.ReconcilerResult.RequeueError(err)
	}
	s.End()

	// == handle delegated zone in Edge DNS
	_, szd := r.Tracer.Start(ctx, "CreateZoneDelegationForExternalDNS")
	err = r.DNSProvider.CreateZoneDelegationForExternalDNS(rs)
	if err != nil {
		log.Err(err).Msg("Unable to create zone delegation")
		m.IncrementError(rs)
		return r.ReconcilerResult.Requeue()
	}
	szd.End()

	// == Status =
	err = r.updateStatus(rs, dnsEndpoint)
	if err != nil {
		m.IncrementError(rs)
		return r.ReconcilerResult.RequeueError(err)
	}
	// == Finish ==========
	// Everything went fine, requeue after some time to catch up
	// with external Gslb status
	m.IncrementReconciliation(rs)
	return r.ReconcilerResult.Requeue()
}
