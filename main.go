package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var err = godotenv.Load()

func main() {
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	fmt.Println(os.Getenv("CLIENT_ID"))
	startServer()
}
