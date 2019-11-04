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

type JWT struct {
	Secret           string
	Duration         int
	RefreshDuration  int
	MaxRefresh       int
	SigningAlgorithm string
}

type Config struct {
	dbConfig    *Database
	bootstraped bool
}

func (c *Config) Bootstrap() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	c.dbConfig = c.loadDatabaseConfig()
	c.bootstraped = true
}

func (c *Config) loadDatabaseConfig() *Database {
	if c.bootstraped {
		c.Bootstrap()
	}
	return &Database{
		PSN:        environment.GetEnv("DATABASE_PSN", "root@(localhost)/test?charset=utf8&parseTime=True&loc=Local"),
		LogQueries: environment.GetEnvAsBool("DATABASE_LOG", false),
		Timeout:    environment.GetEnvAsInt("DATABASE_TIMEOUT", 0),
	}
}

func (c *Config) GetDBConfig() *Database {
	return c.dbConfig
}
