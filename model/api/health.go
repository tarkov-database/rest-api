package api

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/logger"
	"github.com/tarkov-database/rest-api/core/database"
)

var latencyThreshold time.Duration

func init() {
	if env := os.Getenv("UNHEALTHY_LATENCY"); len(env) > 0 {
		d, err := time.ParseDuration(env)
		if err != nil {
			log.Printf("Unhealthy latency value is not valid: %s\n", err)
			os.Exit(2)
		}
		latencyThreshold = d
	} else {
		latencyThreshold = 300 * time.Millisecond
	}
}

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
func GetHealth() *Health {
	svc := &Service{}

	health := &Health{
		OK:      true,
		Service: svc,
	}

	svc.Database = getDatabaseStatus()
	if svc.Database != OK {
		health.OK = false
	}

	return health
}

func getDatabaseStatus() Status {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()

	if err := database.Ping(ctx); err != nil {
		logger.Errorf("Error while checking database connection: %s", err)
		return Failure
	}

	latency := time.Since(start)
	if latency > latencyThreshold {
		logger.Warningf("Database latency exceeds threshold with %s", latency)
		return Warning
	}

	return OK
}
