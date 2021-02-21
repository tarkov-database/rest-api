package health

import (
	"context"
	"os"
	"time"

	"github.com/tarkov-database/rest-api/core/database"

	"github.com/google/logger"
)

var (
	databaseStatus = Failure
)

// GetDBStatus returns the database status of the last health check
func GetDBStatus() Status {
	return databaseStatus
}

func scheduler(t *time.Ticker, c chan os.Signal) {
	for {
		select {
		case <-t.C:
			databaseStatus = getDatabaseStatus()
		case <-c:
			t.Stop()
			return
		}
	}
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
