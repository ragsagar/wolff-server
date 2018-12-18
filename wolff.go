package main

import (
	"github.com/ragsagar/wolff/server"
	"github.com/ragsagar/wolff/store"
)

func main() {
	dataStore := store.NewSQLStore("wolffapp", "password", "wolffapp", "localhost:5432")
	srv := server.NewServer(dataStore)
	srv.Run(":8080")
}
