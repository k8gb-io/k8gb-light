package rs

import (
	"context"
	"reflect"

	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// todo: rename package to reconciliation

type MapperResult int

const (
	MapperResultExists MapperResult = iota
	MapperResultCreate
	MapperResultError
)

type IngressMapper struct {
	c client.Client
}

func NewIngressMapper(c client.Client) *IngressMapper {
	return &IngressMapper{
		c: c,
	}
}

func (i *IngressMapper) Update(state *ReconciliationState) error {
	return i.c.Update(context.TODO(), state.Ingress)
}

func (i *IngressMapper) Get(selector types.NamespacedName) (rs *ReconciliationState, result MapperResult, err error) {
	var ing = &netv1.Ingress{}
	err = i.c.Get(context.TODO(), selector, ing)
	result, err = i.getConverterResult(err)
	if result == MapperResultError {
		return nil, result, err
	}
	rs, err = NewReconciliationState(ing)
	if err != nil {
		result = MapperResultError
	}
	return rs, result, err
}

// Equal compares given ingress annotations and Ingres.Spec. If any of ingresses doesn't exist, returns false
func (i *IngressMapper) Equal(rs1 *ReconciliationState, rs2 *ReconciliationState) bool {
	if rs1 == nil || rs2 == nil {
		return false
	}
	if !reflect.DeepEqual(rs1.Spec, rs2.Spec) {
		return false
	}
	if !reflect.DeepEqual(rs1.Ingress.Spec, rs2.Ingress.Spec) {
		return false
	}
	return true
}

func (i *IngressMapper) getConverterResult(err error) (MapperResult, error) {
	if err != nil && errors.IsNotFound(err) {
		return MapperResultCreate, nil
	} else if err != nil {
		return MapperResultError, err
	}
	return MapperResultExists, nil
}
