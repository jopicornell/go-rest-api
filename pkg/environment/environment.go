package environment

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

var envLoaded bool

func GetEnv(key string, defaultValue string) string {
	if !envLoaded {
		if err := godotenv.Load(); err != nil {
			log.Print("No .env.example file found")
		}
	}
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetEnvAsInt(name string, defaultValue int) int {
	convDefValue := strconv.Itoa(defaultValue)
	if convValue, err := strconv.Atoi(GetEnv(name, convDefValue)); err == nil {
		return convValue
	}
	return defaultValue
}

func GetEnvAsBool(name string, defaultValue bool) bool {
	convDefValue := strconv.FormatBool(defaultValue)
	if convValue, err := strconv.ParseBool(GetEnv(name, convDefValue)); err == nil {
		return convValue
	}
	return defaultValue
}

func GetEnvAsSlice(name string, defaultValue []string) []string {
	convDefValue := strings.Join(defaultValue, "|")
	return strings.Split(GetEnv(name, convDefValue), "|")
}
