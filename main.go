package main

import (
	"fmt"
	"io"

	"github.com/tarkov-database/rest-api/core/database"
	"github.com/tarkov-database/rest-api/core/health"
	"github.com/tarkov-database/rest-api/core/server"
	"github.com/tarkov-database/rest-api/model/api"

	"github.com/google/logger"
)

func main() {
	fmt.Printf("Starting up Tarkov Database REST API %s\n\n", api.Version)

	defLog := logger.Init("default", true, false, io.Discard)
	defer defLog.Close()

	if err := database.Init(); err != nil {
		logger.Fatalf("Database initiation error: %s", err)
	}
	defer func() {
		if err := database.Shutdown(); err != nil {
			logger.Errorf("Database shutdown error: %s", err)
		}
	}()

	health.InitChecks()

	if err := server.ListenAndServe(); err != nil {
		logger.Errorf("HTTP server error: %s", err)
	}
}
