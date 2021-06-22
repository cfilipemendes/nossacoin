package main

import (
	"log"

	"github.com/joho/godotenv"
)

//CHAIN
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//Initialize Chain.
	chainInit()

	log.Fatal(bootstrapServer())
}
