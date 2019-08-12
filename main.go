package main

import (
	"fmt"
	"io/ioutil"

	"github.com/tarkov-database/api/core/database"
	"github.com/tarkov-database/api/core/server"
	"github.com/tarkov-database/api/model/api"

	"github.com/google/logger"
)

func main() {
	fmt.Printf("Starting up Tarkov Database API %s\n\n", api.Version)

	defLog := logger.Init("default", true, false, ioutil.Discard)
	defer defLog.Close()

	database.Init()
	server.ListenAndServe()
}
