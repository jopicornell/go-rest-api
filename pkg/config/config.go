package config

import (
	"github.com/joho/godotenv"
	"github.com/jopicornell/go-rest-api/pkg/environment"
	"log"
)

type Database struct {
	PSN        string
	LogQueries bool
	Timeout    int
}

type Server struct {
	Port         string
	Debug        bool
	ReadTimeout  int
	WriteTimeout int
}

type JWT struct {
	Secret           string
	Duration         int
	RefreshDuration  int
	MaxRefresh       int
	SigningAlgorithm string
}

var dbConfig *Database = loadDatabaseConfig()

var environmentFileParsed = false

func load() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	environmentFileParsed = true
}

func loadDatabaseConfig() *Database {
	if environmentFileParsed {
		load()
	}
	return &Database{
		PSN:        environment.GetEnv("DATABASE_PSN", "root@(localhost)/test?charset=utf8&parseTime=True&loc=Local"),
		LogQueries: environment.GetEnvAsBool("DATABASE_LOG", false),
		Timeout:    environment.GetEnvAsInt("DATABASE_TIMEOUT", 0),
	}
}

func GetDBConfig() *Database {
	return dbConfig
}
