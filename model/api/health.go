package api

import (
	"context"
	"time"

	"github.com/tarkov-database/rest-api/core/database"
)

// Status represents the status code of a service
type Status int

const (
	// OK status if all checks were successful
	OK Status = iota

	// Warning status if non-critical issues are discovered
	Warning

	// Failure status when critical problems are discovered
	Failure
)

// Health represents the object of the health root endpoint
type Health struct {
	OK      bool     `json:"ok"`
	Service *Service `json:"service"`
}

// Service holds all services with their respective status
type Service struct {
	Database Status `json:"database"`
}

// GetHealth performs a self-check and returns the result
func GetHealth() (*Health, error) {
	var err error
	var ok = true

	svc := &Service{}

	svc.Database, err = getDatabaseStatus()
	if err != nil {
		return nil, err
	}

	if svc.Database != OK {
		ok = false
	}

	health := &Health{
		OK:      ok,
		Service: svc,
	}

	return health, nil
}

func getDatabaseStatus() (Status, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()

	if err := database.Ping(ctx); err != nil {
		return Failure, err
	}

	if time.Since(start).Seconds() > 3 {
		return Warning, nil
	}

	return OK, nil
}
