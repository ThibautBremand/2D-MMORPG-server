package main

import (
	"fmt"
	"log"
	"server/db"
	"server/webserver"

	"github.com/joho/godotenv"
)

// startDatabases connects to the storage and the redis.
func startDatabases() {
	if err := db.Open(); err != nil {
		fmt.Printf("error %v", err)
	}
	db.Start()
}

// init is invoked before main(), and loads values from .env as env variables.
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	webserver.PrepareConfig()
	startDatabases()
	defer db.Close()

	webserver.StartHub()
	webserver.StartWebServer()
}
