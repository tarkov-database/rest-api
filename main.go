package main

import (
	"fmt"
	"io/ioutil"

	"github.com/tarkov-database/rest-api/core/database"
	"github.com/tarkov-database/rest-api/core/server"
	"github.com/tarkov-database/rest-api/model/api"

	"github.com/google/logger"
)

func main() {
	fmt.Printf("Starting up Tarkov Database REST API %s\n\n", api.Version)

	defLog := logger.Init("default", true, false, ioutil.Discard)
	defer defLog.Close()

	err := database.Init()
	if err != nil {
		logger.Fatal("Database initiation error: %s", err)
	}

	err = server.ListenAndServe()
	if err != nil {
		logger.Fatal("HTTP server error: %s", err)
	}
}
