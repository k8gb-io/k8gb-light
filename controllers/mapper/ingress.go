package mapper

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
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	externaldns "sigs.k8s.io/external-dns/endpoint"

	"cloud.example.com/annotation-operator/controllers/depresolver"

	"cloud.example.com/annotation-operator/controllers/providers/metrics"
	corev1 "k8s.io/api/core/v1"

	"cloud.example.com/annotation-operator/controllers/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IngressMapper provides API for working with ingress
type IngressMapper struct {
	c      client.Client
	config *depresolver.Config
	rs     *LoopState
}

func NewIngressMapper(c client.Client, config *depresolver.Config) *IngressMapper {
	return &IngressMapper{
		c:      c,
		config: config,
	}
}

func (i *IngressMapper) UpdateStatusAnnotation() (err error) {
	// check if object has not been deleted
	var r Result
	var s *LoopState
	s, r, err = NewCommonProvider(i.c, i.config).Get(i.rs.NamespacedName)
	switch r {
	case ResultError:
		return err
	case ResultNotFound:
		// object was deleted
		return nil
	}
	i.rs.Status = i.GetStatus()
	// don't do update if nothing has changed
	if s.Ingress.Annotations[AnnotationStatus] == i.rs.Status.String() {
		return nil
	}
	// update the planned object
	s.Ingress.Annotations[AnnotationStatus] = i.rs.Status.String()
	return i.c.Update(context.TODO(), s.Ingress)
}

// Equal compares given ingress annotations and Ingres.Spec. If any of ingresses doesn't exist, returns false
func (i *IngressMapper) Equal(rs *LoopState) bool {
	if i.rs == nil || rs == nil {
		return false
	}
	if !reflect.DeepEqual(i.rs.Spec, rs.Spec) {
		return false
	}
	if !reflect.DeepEqual(i.rs.Ingress.Spec, rs.Ingress.Spec) {
		return false
	}
	return true
}

func (i *IngressMapper) TryInjectFinalizer() (Result, error) {
	if i.rs == nil || i.rs.Ingress == nil {
		return ResultError, fmt.Errorf("injecting finalizer from nil values")
	}
	if !utils.Contains(i.rs.Ingress.GetFinalizers(), Finalizer) {
		i.rs.Ingress.SetFinalizers(append(i.rs.Ingress.GetFinalizers(), Finalizer))
		err := i.c.Update(context.TODO(), i.rs.Ingress)
		if err != nil {
			return ResultError, err
		}
		return ResultFinalizerInstalled, nil
	}
	return ResultContinue, nil
}

func (i *IngressMapper) TryRemoveFinalizer(finalize func(*LoopState) error) (Result, error) {
	if i.rs == nil || i.rs.Ingress == nil {
		return ResultError, fmt.Errorf("removing finalizer from nil values")
	}
	if utils.Contains(i.rs.Ingress.GetFinalizers(), Finalizer) {
		isMarkedToBeDeleted := i.rs.Ingress.GetDeletionTimestamp() != nil
		if !isMarkedToBeDeleted {
			return ResultContinue, nil
		}
		err := finalize(i.rs)
		if err != nil {
			return ResultError, err
		}
		i.rs.Ingress.SetFinalizers(utils.Remove(i.rs.Ingress.GetFinalizers(), Finalizer))
		err = i.c.Update(context.TODO(), i.rs.Ingress)
		if err != nil {
			return ResultError, err
		}
		return ResultFinalizerRemoved, nil
	}
	return ResultContinue, nil
}

func (i *IngressMapper) getHealthStatus() map[string]metrics.HealthStatus {
	serviceHealth := make(map[string]metrics.HealthStatus)
	for _, rule := range i.rs.Ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			if path.Backend.Service == nil || path.Backend.Service.Name == "" {
				serviceHealth[rule.Host] = metrics.NotFound
				continue
			}

			// check if service exists
			selector := types.NamespacedName{Namespace: i.rs.NamespacedName.Namespace, Name: path.Backend.Service.Name}
			service := &corev1.Service{}
			err := i.c.Get(context.TODO(), selector, service)
			if err != nil {
				if errors.IsNotFound(err) {
					serviceHealth[rule.Host] = metrics.NotFound
					continue
				}
				return make(map[string]metrics.HealthStatus, 0)
			}

			// check if service endpoint exists
			ep := &corev1.Endpoints{}
			err = i.c.Get(context.TODO(), selector, ep)
			if err != nil {
				return serviceHealth
			}
			serviceHealth[rule.Host] = metrics.Unhealthy
			for _, subset := range ep.Subsets {
				if len(subset.Addresses) > 0 {
					serviceHealth[rule.Host] = metrics.Healthy
				}
			}
		}
	}
	return serviceHealth
}

func (i *IngressMapper) GetExposedIPs() ([]string, error) {
	var exposed []string
	for _, ing := range i.rs.Ingress.Status.LoadBalancer.Ingress {
		if len(ing.IP) > 0 {
			exposed = append(exposed, ing.IP)
		}
		if len(ing.Hostname) > 0 {
			ips, err := utils.Dig(ing.Hostname, i.config.EdgeDNSServers...)
			if err != nil {
				return nil, err
			}
			exposed = append(exposed, ips...)
		}
	}
	return exposed, nil
}

func (i *IngressMapper) getHealthyRecords() map[string][]string {
	// TODO: make mapper for DNSEndpoint
	healthyRecords := make(map[string][]string)
	dnsEndpoint := &externaldns.DNSEndpoint{}
	err := i.c.Get(context.TODO(), i.rs.NamespacedName, dnsEndpoint)
	if err != nil {
		return healthyRecords
	}

	serviceRegex := regexp.MustCompile("^localtargets")
	for _, endpoint := range dnsEndpoint.Spec.Endpoints {
		local := serviceRegex.Match([]byte(endpoint.DNSName))
		if !local && endpoint.RecordType == RecordTypeA {
			if len(endpoint.Targets) > 0 {
				healthyRecords[endpoint.DNSName] = endpoint.Targets
			}
		}
	}
	return healthyRecords
}

func (i *IngressMapper) GetStatus() (status Status) {
	csv := func(rs *LoopState) string {
		var hosts []string
		for _, r := range rs.Ingress.Spec.Rules {
			hosts = append(hosts, r.Host)
		}
		return strings.Join(hosts, ", ")
	}

	return Status{
		ServiceHealth:  i.getHealthStatus(),
		HealthyRecords: i.getHealthyRecords(),
		GeoTag:         i.config.ClusterGeoTag,
		Hosts:          csv(i.rs),
	}
}

func (i *IngressMapper) SetReference(rs *LoopState) {
	i.rs = rs
}
