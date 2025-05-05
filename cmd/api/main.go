package main

import (
	"fmt"
	"os"

	"github.com/0xirvan/go-jwt-auth/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Fatalln("Error loading .env file")
	}

	fmt.Println(os.Getenv("foo"))
}
