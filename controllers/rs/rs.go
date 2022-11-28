package rs

import (
	"fmt"
	"strconv"

	"k8s.io/apimachinery/pkg/types"

	netv1 "k8s.io/api/networking/v1"

	"cloud.example.com/annotation-operator/controllers/depresolver"
	"k8s.io/apimachinery/pkg/util/json"
)

const (
	PrimaryGeoTagAnnotation              = "k8gb.io/primary-geotag"
	StrategyAnnotation                   = "k8gb.io/strategy"
	DnsTTLSecondsAnnotation              = "k8gb.io/dns-ttl-seconds"
	SplitBrainThresholdSecondsAnnotation = "k8gb.io/splitbrain-threshold-seconds"
	WeightAnnotationJSON                 = "k8gb.io/weights"
)

type Spec struct {
	PrimaryGeoTag              string         `json:"primaryGeoTag"`
	Type                       string         `json:"strategy"`
	DNSTtlSeconds              int            `json:"dnsTTLSeconds"`
	SplitBrainThresholdSeconds int            `json:"splitBrainThresholdSeconds"`
	Weights                    map[string]int `json:"weights"`
}

func (s *Spec) String() string {
	return fmt.Sprintf("strategy: %s, geo: %s", s.Type, s.PrimaryGeoTag)
}

type ReconciliationState struct {
	Ingress        *netv1.Ingress
	Spec           Spec
	NamespacedName types.NamespacedName
}

func NewReconciliationState(ingress *netv1.Ingress) (m *ReconciliationState, err error) {
	m = new(ReconciliationState)
	if ingress == nil {
		return m, fmt.Errorf("nil *ingress")
	}
	m.Ingress = ingress
	m.Spec, err = m.asSpec(ingress.GetAnnotations())
	m.NamespacedName = types.NamespacedName{Namespace: ingress.Namespace, Name: ingress.Name}
	return m, err
}

func (rs *ReconciliationState) HasStrategy() bool {
	annotations := rs.Ingress.GetAnnotations()
	_, found := annotations[StrategyAnnotation]
	return found
}

func (rs *ReconciliationState) asAnnotation(s Spec) (annotations map[string]string, err error) {
	var predefinedStrategy = Spec{
		DNSTtlSeconds:              30,
		SplitBrainThresholdSeconds: 300,
	}
	annotations = make(map[string]string, 0)
	if s.DNSTtlSeconds != predefinedStrategy.DNSTtlSeconds {
		annotations[DnsTTLSecondsAnnotation] = strconv.Itoa(s.DNSTtlSeconds)
	}
	if s.SplitBrainThresholdSeconds != predefinedStrategy.SplitBrainThresholdSeconds {
		annotations[SplitBrainThresholdSecondsAnnotation] = strconv.Itoa(s.SplitBrainThresholdSeconds)
	}
	annotations[PrimaryGeoTagAnnotation] = s.PrimaryGeoTag
	annotations[StrategyAnnotation] = s.Type

	weights, err := json.Marshal(s.Weights)
	if err != nil {
		return annotations, fmt.Errorf("reading weights %v", err)
	}
	annotations[WeightAnnotationJSON] = string(weights)
	return annotations, err
}

func (rs *ReconciliationState) asSpec(annotations map[string]string) (result Spec, err error) {
	toInt := func(k string, v string) (int, error) {
		intValue, err := strconv.Atoi(v)
		if err != nil {
			return -1, fmt.Errorf("can't parse annotation value %s to int for key %s", v, k)
		}
		return intValue, nil
	}
	result = Spec{
		Type:                       "",
		PrimaryGeoTag:              "",
		DNSTtlSeconds:              30,
		SplitBrainThresholdSeconds: 300,
		Weights:                    nil,
	}

	if value, found := annotations[StrategyAnnotation]; found {
		result.Type = value
	}
	if value, found := annotations[PrimaryGeoTagAnnotation]; found {
		result.PrimaryGeoTag = value
	}
	if value, found := annotations[SplitBrainThresholdSecondsAnnotation]; found {
		if result.SplitBrainThresholdSeconds, err = toInt(SplitBrainThresholdSecondsAnnotation, value); err != nil {
			return result, err
		}
	}
	if value, found := annotations[DnsTTLSecondsAnnotation]; found {
		if result.DNSTtlSeconds, err = toInt(DnsTTLSecondsAnnotation, value); err != nil {
			return result, err
		}
	}

	if value, found := annotations[WeightAnnotationJSON]; found {
		w := make(map[string]int, 0)
		err = json.Unmarshal([]byte(value), &w)
		if err != nil {
			return result, err
		}
		result.Weights = w
	}

	if result.Type == depresolver.FailoverStrategy {
		if len(result.PrimaryGeoTag) == 0 {
			return result, fmt.Errorf("%s strategy requires annotation %s", depresolver.FailoverStrategy, PrimaryGeoTagAnnotation)
		}
	}
	return result, nil
}

func (rs *ReconciliationState) Save() error {
	rs.
}
