package main

import (
	"fmt"

	"github.com/ragsagar/wolff/server"
	"github.com/ragsagar/wolff/store"
)

func main() {
	fmt.Println("Hello World!")
	dataStore := store.NewSQLStore("wolffapp", "password", "wolffapp", "localhost:5432")
	srv := server.NewServer(dataStore)
	srv.Run(":8080")
}
