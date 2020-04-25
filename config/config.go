package config

import (
	"log"
	"os"
)

func EnvVar(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("%s env var required!", key)
	}

	return value
}