package converter

import (
	"context"

	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type V1Beta1 struct {
	c client.Client
}

func NewIngressV1Beta1(c client.Client) *V1Beta1 {
	return &V1Beta1{
		c: c,
	}
}

func (l *V1Beta1) ListIngresses(a client.Object) (nn []types.NamespacedName, err error) {
	ingList := &netv1beta1.IngressList{}
	err = l.c.List(context.TODO(), ingList, client.InNamespace(a.GetNamespace()))
	if err != nil {
		return nn, err
	}
	for _, ing := range ingList.Items {
		for _, rule := range ing.Spec.Rules {
			for _, path := range rule.HTTP.Paths {
				if path.Backend.ServiceName == a.GetName() {
					nn = append(nn, types.NamespacedName{Namespace: a.GetNamespace(), Name: ing.Name})
				}
			}
		}
	}
	return nn, err
}
