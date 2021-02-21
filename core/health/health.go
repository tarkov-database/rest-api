package health

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	updateInterval   time.Duration
	latencyThreshold time.Duration
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

func init() {
	if env := os.Getenv("HEALTHCHECK_INTERVAL"); len(env) > 0 {
		d, err := time.ParseDuration(env)
		if err != nil {
			log.Printf("Health check interval value is not valid: %s\n", err)
			os.Exit(2)
		}
		updateInterval = d
	} else {
		updateInterval = 30 * time.Second
	}

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

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(updateInterval)

	go scheduler(ticker, sig)
}
