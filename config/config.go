package config

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadEnvs(filename ...string) {
	if filename == nil {
		filename = []string{"./.env"}
	}
	if os.Getenv("APP") == "" {
		err := godotenv.Load(filename[0])
		if err != nil {
			panic("Error loading .env file")
		}
	}
}

func IsLocal() bool {
	return os.Getenv("APP") == "dev"
}
