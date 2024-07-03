package main

import (
	"log"
	"os"

	"go_final_project/pkg/handlers"

	"github.com/joho/godotenv"
)

const(
	envFile = ".env"
)

func main() {
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error load env: %s", err)
	}
	port := os.Getenv("TODO_PORT")
	dbPath := os.Getenv("TODO_DBFILE")
	
	handlers.HandleRequests(dbPath, port)
}
