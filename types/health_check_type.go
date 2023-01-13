package types

// Error is a type for error
type HealthCheck struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}
