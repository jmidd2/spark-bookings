package config

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnvString(key string) string {
	return os.Getenv(key)
}

func GetEnvStringOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetEnvInt(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return 0
	}
	return value
}

func GetEnvBool(key string) bool {
	value, err := strconv.ParseBool(os.Getenv(key))
	if err != nil {
		return false
	}
	return value
}

func GetEnvStringRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Environment variable %s is required", key))
	}
	return value
}

func GetEnvIntRequired(key string) int {
	value := GetEnvInt(key)
	if value == 0 {
		panic(fmt.Sprintf("Environment variable %s is required", key))
	}
	return value
}

func GetEnvBoolRequired(key string) bool {
	value := GetEnvBool(key)
	if !value {
		panic(fmt.Sprintf("Environment variable %s is required", key))
	}
	return value
}
