package converter

import (
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IIngress interface {
	ListIngresses(object client.Object) ([]types.NamespacedName, error)
}
