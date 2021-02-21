package health

import (
	"context"
	"os"
	"time"

	"github.com/tarkov-database/rest-api/core/database"

	"github.com/google/logger"
)

var (
	databaseStatus = OK
)

// DatabaseStatus returns the database status of the last health check
func DatabaseStatus() Status {
	return databaseStatus
}

func scheduler(t *time.Ticker, c chan os.Signal) {
	for {
		select {
		case <-t.C:
			updateStatus()
		case <-c:
			t.Stop()
			return
		}
	}
}

func updateStatus() {
	databaseStatus = getDatabaseStatus()
}

func getDatabaseStatus() Status {
	timeout := 30 * time.Second
	if cfg.updateInterval < timeout {
		timeout = cfg.updateInterval
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()

	if err := database.Ping(ctx); err != nil {
		logger.Errorf("Error while checking database connection: %s", err)
		return Failure
	}

	latency := time.Since(start)
	if latency > cfg.latencyThreshold {
		logger.Warningf("Database latency exceeds threshold with %s", latency)
		return Warning
	}

	return OK
}
