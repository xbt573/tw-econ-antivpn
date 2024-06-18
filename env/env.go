package env

import (
	"log"
	"os"
)

func Get(name string) string {
	env, exists := os.LookupEnv(name)
	if !exists {
		log.Fatalf("%v not set\n", name)
	}

	return env
}

func GetDefault(name string, defaultValue string) string {
	env, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}

	return env
}
