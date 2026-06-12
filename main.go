package main

import (
	"log"

	"github.com/joho/godotenv"
	"gitsearch/cmd"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	cmd.Execute()
}
