package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ragsagar/wolff/server"
	"github.com/ragsagar/wolff/store"
	"github.com/spf13/viper"
)

func main() {
	viper.AddConfigPath("config")
	env := os.Getenv("WOLFF_ENV")
	if env == "PRODUCTION" {
		viper.SetConfigName("production")
	} else if env == "STAGING" {
		viper.SetConfigName("staging")
	} else {
		viper.SetConfigName("dev")
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Fatal error reading config file", err)
	}
	dataStore := store.NewSQLStore(viper.GetString("DB_NAME"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_USER"), viper.GetString("DB_SERVER"))
	srv := server.NewServer(dataStore)
	srv.Run(fmt.Sprintf(":%s", viper.GetString("APP_SERVER_PORT")))
}
