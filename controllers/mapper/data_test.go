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
	"fmt"

	"github.com/k8gb-io/k8gb-light/controllers/depresolver"

	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	externaldns "sigs.k8s.io/external-dns/endpoint"
)

var Prefix = netv1.PathTypePrefix

type Data struct {
	Ingress                 *netv1.Ingress
	LocalTargetsDNSEndpoint *externaldns.DNSEndpoint
	Service                 *corev1.Service
	Endpoint                *corev1.Endpoints
}

var testData = Data{
	Ingress: &netv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{Name: "ing"},
		Spec: netv1.IngressSpec{
			Rules: []netv1.IngressRule{
				{
					Host: "demo.cloud.example.com",
					IngressRuleValue: netv1.IngressRuleValue{
						HTTP: &netv1.HTTPIngressRuleValue{
							Paths: []netv1.HTTPIngressPath{
								{Path: "/", PathType: &Prefix, Backend: netv1.IngressBackend{
									Service: &netv1.IngressServiceBackend{
										Name: "frontend-podinfo",
										Port: netv1.ServiceBackendPort{
											Name:   "http",
											Number: 9898,
										},
									},
								}},
							},
						},
					},
				},
			}},
	},

	LocalTargetsDNSEndpoint: &externaldns.DNSEndpoint{
		ObjectMeta: metav1.ObjectMeta{Name: "ing",
			OwnerReferences: []metav1.OwnerReference{{Name: "ing", Kind: "Ingress"}}},
		Spec: externaldns.DNSEndpointSpec{
			Endpoints: []*externaldns.Endpoint{
				{Targets: []string{"172.18.0.5", "172.18.0.6"}, DNSName: "localtargets-demo.cloud.example.com", RecordType: "A"},
				{Targets: []string{"172.18.0.5", "172.18.0.6", "172.18.0.3", "172.18.0.4"},
					DNSName: "demo.cloud.example.com", RecordType: "A", Labels: map[string]string{"strategy": depresolver.RoundRobinStrategy}},
			},
		},
	},

	Service: &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "frontend-podinfo"},
		Spec: corev1.ServiceSpec{
			Ports:       []corev1.ServicePort{{Name: "http", Port: 9898, Protocol: corev1.ProtocolTCP}},
			ClusterIP:   "10.43.174.228",
			ClusterIPs:  []string{"10.43.174.228"},
			Type:        "ClusterIP",
			ExternalIPs: nil,
		},
	},

	Endpoint: &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{Name: "frontend-podinfo"},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{
					{
						IP:        "10.42.0.9",
						TargetRef: &corev1.ObjectReference{Kind: "Pod"},
					},
				},
			},
		},
	},
}

func (d Data) deepCopy() Data {
	return Data{
		Ingress:                 d.Ingress.DeepCopy(),
		Service:                 d.Service.DeepCopy(),
		Endpoint:                d.Endpoint.DeepCopy(),
		LocalTargetsDNSEndpoint: d.LocalTargetsDNSEndpoint.DeepCopy(),
	}
}

func FOon2c1() Data {
	data := testData.deepCopy()
	data.Ingress.Annotations =
		map[string]string{AnnotationStrategy: depresolver.FailoverStrategy, AnnotationPrimaryGeoTag: "eu"}
	data.LocalTargetsDNSEndpoint.Spec.Endpoints[0].Targets = []string{"172.18.0.5", "172.18.0.6"}
	data.LocalTargetsDNSEndpoint.Spec.Endpoints[1].Targets = []string{"172.18.0.3", "172.18.0.4"}
	return data
}

func FOon2c2() Data {
	data := testData.deepCopy()
	data.Ingress.Annotations =
		map[string]string{AnnotationStrategy: depresolver.FailoverStrategy, AnnotationPrimaryGeoTag: "eu"}
	data.LocalTargetsDNSEndpoint.Spec.Endpoints[0].Targets = []string{"172.18.0.3", "172.18.0.4"}
	data.LocalTargetsDNSEndpoint.Spec.Endpoints[1].Targets = []string{"172.18.0.3", "172.18.0.4"}
	return data
}

// RoundRobin on two clusters
func RRon2() Data {
	data := testData.deepCopy()
	data.Ingress.Annotations = map[string]string{AnnotationStrategy: depresolver.RoundRobinStrategy}
	data.Ingress.Status.LoadBalancer.Ingress = []corev1.LoadBalancerIngress{{IP: "172.18.0.5"}, {IP: "172.18.0.6"}}
	return data
}

func (d Data) AddHost(host string, localtargets []string, targets []string) Data {
	rule := netv1.IngressRule{
		Host: host,
		IngressRuleValue: netv1.IngressRuleValue{
			HTTP: &netv1.HTTPIngressRuleValue{
				Paths: []netv1.HTTPIngressPath{
					{Path: "/", PathType: &Prefix, Backend: netv1.IngressBackend{
						Service: &netv1.IngressServiceBackend{
							Name: "frontend-podinfo",
							Port: netv1.ServiceBackendPort{
								Name:   "http",
								Number: 9898,
							},
						},
					}},
				},
			},
		},
	}
	eps := []*externaldns.Endpoint{
		{
			DNSName:    fmt.Sprintf("localtargets-%s", host),
			RecordType: RecordTypeA,
			Targets:    targets,
		},
		{
			DNSName:    host,
			RecordType: RecordTypeA,
			Targets:    targets,
			Labels:     map[string]string{"strategy": depresolver.RoundRobinStrategy},
		},
	}
	d1 := d.deepCopy()
	d1.Ingress.Spec.Rules = append(d1.Ingress.Spec.Rules, rule)
	d1.LocalTargetsDNSEndpoint.Spec.Endpoints = append(d1.LocalTargetsDNSEndpoint.Spec.Endpoints, eps...)
	return d1
}
