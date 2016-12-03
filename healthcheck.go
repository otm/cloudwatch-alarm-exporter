package main

// HealthChecker is provides an interface for healthchecking
type HealthChecker interface {
	Health() error
}

// AlwaysHealthy is a healthchecker that will never fail
type AlwaysHealthy struct{}

// Health always return a nil error
func (h AlwaysHealthy) Health() error {
	return nil
}
