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

type Server struct {
	ApiUrl      string
	StaticsPath string
	Port        string
	JWTSecret   string
}

type Config struct {
	dbConfig     *Database
	serverConfig *Server
	bootstraped  bool
}

func (c *Config) Bootstrap() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	c.dbConfig = c.loadDatabaseConfig()
	c.serverConfig = c.loadServerConfig()
	c.bootstraped = true
}

func (c *Config) loadDatabaseConfig() *Database {
	return &Database{
		PSN:        environment.GetEnv("DATABASE_PSN", "root@(localhost)/test?charset=utf8&parseTime=True&loc=Local"),
		LogQueries: environment.GetEnvAsBool("DATABASE_LOG", false),
		Timeout:    environment.GetEnvAsInt("DATABASE_TIMEOUT", 0),
	}
}

func (c *Config) loadServerConfig() *Server {
	return &Server{
		Port:        environment.GetEnv("PORT", "8080"),
		ApiUrl:      environment.GetEnv("API_URL", "/api"),
		StaticsPath: environment.GetEnv("STATICS_PATH", "static"),
		JWTSecret:   environment.GetEnv("JWT_SECRET", "secret"),
	}
}

func (c *Config) GetDBConfig() *Database {
	return c.dbConfig
}

func (c *Config) GetServerConfig() *Server {
	return c.serverConfig
}
