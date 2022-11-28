package status

type HealthStatus string

const (
	Healthy   HealthStatus = "Healthy"
	Unhealthy HealthStatus = "Unhealthy"
	NotFound  HealthStatus = "NotFound"
)

func (h HealthStatus) String() string {
	return string(h)
}
