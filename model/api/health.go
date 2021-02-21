package api

import (
	"github.com/tarkov-database/rest-api/core/health"
)

// Health represents the object of the health root endpoint
type Health struct {
	OK      bool     `json:"ok"`
	Service *Service `json:"service"`
}

// Service holds all services with their respective status
type Service struct {
	Database health.Status `json:"database"`
}

// GetHealth performs a self-check and returns the result
func GetHealth() *Health {
	svc := &Service{}

	h := &Health{
		OK:      true,
		Service: svc,
	}

	svc.Database = health.GetDBStatus()
	if svc.Database != health.OK {
		h.OK = false
	}

	return h
}
