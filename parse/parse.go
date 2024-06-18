package parse

import (
	"log"
	"strconv"
)

func GetIntOrFail(str string) int {
	result, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalf("\"%v\" is invalid integer\n", str)
	}

	return result
}
