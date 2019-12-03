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

type Redis struct {
	Host     string
	Password string
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
	redisConfig  *Redis
	serverConfig *Server
	bootstraped  bool
}

func (c *Config) Bootstrap() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env.example file found")
	}
	c.dbConfig = c.loadDatabaseConfig()
	c.redisConfig = c.loadRedisConfig()
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

func (c *Config) loadRedisConfig() *Redis {
	return &Redis{
		Host:     environment.GetEnv("REDIS_HOST", "localhost:6379"),
		Password: environment.GetEnv("REDIS_PASSWORD", "false"),
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

func (c *Config) GetRedisConfig() *Redis {
	return c.redisConfig
}

func (c *Config) GetServerConfig() *Server {
	return c.serverConfig
}
