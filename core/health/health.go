package health

import (
	"os"
	"os/signal"
	"syscall"
	"time"
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

// InitChecks initiates health check jobs
func InitChecks() {
	updateStatus()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(cfg.updateInterval)

	go scheduler(ticker, sig)
}
