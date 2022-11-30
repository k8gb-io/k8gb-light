package reconciliation

import "k8s.io/apimachinery/pkg/types"

// Mapper is wrapper around resource. Mappers are an only way to access resources
type Mapper interface {
	UpdateStatus(state *ReconciliationState) error
	Get(selector types.NamespacedName) (rs *ReconciliationState, result MapperResult, err error)
	Equal(rs1 *ReconciliationState, rs2 *ReconciliationState) bool
}
