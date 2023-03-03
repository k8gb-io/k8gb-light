package converter

import (
	"context"

	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type V1 struct {
	c client.Client
}

func NewIngressV1(c client.Client) *V1 {
	return &V1{
		c: c,
	}
}

func (l *V1) ListIngresses(a client.Object) (nn []types.NamespacedName, err error) {
	ingList := &netv1.IngressList{}
	err = l.c.List(context.TODO(), ingList, client.InNamespace(a.GetNamespace()))
	if err != nil {
		return nn, err
	}
	for _, ing := range ingList.Items {
		for _, rule := range ing.Spec.Rules {
			for _, path := range rule.HTTP.Paths {
				if path.Backend.Service != nil && path.Backend.Service.Name == a.GetName() {
					nn = append(nn, types.NamespacedName{Namespace: a.GetNamespace(), Name: ing.Name})
				}
			}
		}
	}
	return nn, err
}
