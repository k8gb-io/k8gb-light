package converter

import (
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IIngress interface {
	ListIngresses(object client.Object) ([]types.NamespacedName, error)
}

func EnsureV1(a client.Object) *netv1.Ingress {
	version := a.GetResourceVersion()
	switch version {
	case "networking.k8s.io/v1":
		return a.(*netv1.Ingress)
	case "extensions/v1beta1", "networking.k8s.io/v1beta1":
		return convert(a)
	}
	return a.(*netv1.Ingress)
}

func convert(a client.Object) *netv1.Ingress {

	convertBackend := func(v1beta1 *netv1beta1.IngressBackend) *netv1.IngressBackend {
		if v1beta1 == nil {
			return nil
		}
		return &netv1.IngressBackend{
			Resource: v1beta1.Resource,
			Service: &netv1.IngressServiceBackend{
				Name: v1beta1.ServiceName,
				Port: netv1.ServiceBackendPort{
					Name:   "",
					Number: v1beta1.ServicePort.IntVal,
				},
			},
		}
	}

	v1 := &netv1.Ingress{}
	v1beta1 := a.(*netv1beta1.Ingress)
	v1.Annotations = v1beta1.Annotations
	v1.CreationTimestamp = v1beta1.CreationTimestamp
	v1.DeletionGracePeriodSeconds = v1beta1.DeletionGracePeriodSeconds
	v1.DeletionTimestamp = v1beta1.DeletionTimestamp
	v1.Finalizers = v1beta1.Finalizers
	v1.Generation = v1beta1.Generation
	v1.GenerateName = v1beta1.GenerateName
	v1.Kind = v1beta1.Kind
	v1.Labels = v1beta1.Labels
	v1.ManagedFields = v1beta1.ManagedFields
	v1.Name = v1beta1.Name
	v1.Namespace = v1beta1.Namespace
	v1.ObjectMeta = v1beta1.ObjectMeta
	v1.OwnerReferences = v1beta1.OwnerReferences
	v1.ResourceVersion = v1beta1.ResourceVersion
	v1.TypeMeta = v1beta1.TypeMeta
	v1.UID = v1beta1.UID

	// v1.Spec
	v1.Spec = netv1.IngressSpec{
		IngressClassName: v1beta1.Spec.IngressClassName,

		TLS:            []netv1.IngressTLS{},
		Rules:          []netv1.IngressRule{},
		DefaultBackend: convertBackend(v1beta1.Spec.Backend),
	}

	for _, v := range v1beta1.Spec.TLS {
		tls := netv1.IngressTLS{
			Hosts:      v.Hosts,
			SecretName: v.SecretName,
		}
		v1.Spec.TLS = append(v1.Spec.TLS, tls)
	}

	for _, v := range v1beta1.Spec.Rules {
		rule := netv1.IngressRule{
			Host:             v.Host,
			IngressRuleValue: netv1.IngressRuleValue{},
		}

		if v.HTTP != nil {
			rule.HTTP = &netv1.HTTPIngressRuleValue{}
			for _, p := range v.HTTP.Paths {
				path := netv1.HTTPIngressPath{
					Path:     p.Path,
					PathType: (*netv1.PathType)(p.PathType),
					Backend:  *convertBackend(&p.Backend),
				}
				rule.HTTP.Paths = append(rule.HTTP.Paths, path)
			}
		}
		v1.Spec.Rules = append(v1.Spec.Rules, rule)
	}

	v1.Status = netv1.IngressStatus{
		LoadBalancer: netv1.IngressLoadBalancerStatus{},
	}

	for _, ilv := range v1beta1.Status.LoadBalancer.Ingress {
		is := netv1.IngressLoadBalancerIngress{
			IP:       ilv.IP,
			Hostname: ilv.Hostname,
		}
		for _, ilp := range ilv.Ports {
			isp := netv1.IngressPortStatus{
				Port:     ilp.Port,
				Protocol: ilp.Protocol,
				Error:    ilp.Error,
			}
			is.Ports = append(is.Ports, isp)
		}
		v1.Status.LoadBalancer.Ingress = append(v1.Status.LoadBalancer.Ingress, is)
	}

	return v1
}
