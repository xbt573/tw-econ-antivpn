package env

import (
	"log"
	"os"

	"github.com/xbt573/tw-econ-antivpn/parse"
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

func GetIntDefault(name string, defaultValue int) int {
	env, exists := os.LookupEnv(name)
	if !exists {
		return defaultValue
	}

	return parse.GetIntOrFail(env)
}
