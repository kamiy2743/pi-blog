package config

import (
	"log"
	"os"
)

func mustGetEnvString(envName string) string {
	value := os.Getenv(envName)
	if value == "" {
		log.Fatalf(".env に %s が未設定です。", envName)
	}
	return value
}
