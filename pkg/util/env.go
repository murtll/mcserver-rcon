package util

import (
	"log"
	"os"
	"strconv"
)

func GetStrOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	} else {
		return defaultValue
	}
}

func GetIntOrDefault(key string, defaultValue int) int {
	if stringValue, ok := os.LookupEnv(key); ok {
		value, err := strconv.Atoi(stringValue)
		if err != nil {
			log.Printf("Error converting string env '%s' to int. Falling back to default.\n", key)
			return defaultValue
		}
		return value
	} else {
		return defaultValue
	}
}
