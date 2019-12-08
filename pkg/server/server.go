package server

import (
	goContext "context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/pkg/config"
	"github.com/jopicornell/go-rest-api/pkg/database"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"runtime/debug"
	"time"
)

type Server interface {
	http.Handler
	database.ManagesDatabases
	Close()
	GetServerConfig() *config.Server
	GetRouter() Router
	AddHandler(handler Handler)
	AddStatics(path string, staticPath string)
	ListenAndServe()
}

type server struct {
	http.Server
	config.Config
	relationalDB *database.Postgres
	cache        database.Cache
	Router       Router
}

func Initialize() Server {
	server := &server{}
	server.Config.Bootstrap()
	server.initializeDatabases()
	server.Router = NewRouter(server)
	server.Server = http.Server{
		Addr:    ":" + server.GetServerConfig().Port,
		Handler: server,
	}
	return server
}

func (s *server) GetRelationalDatabase() *sqlx.DB {
	return s.relationalDB.GetDB()
}

func (s *server) GetCache() database.Cache {
	return s.cache
}

func (s *server) AddStatics(exposePath string, staticPath string) {
	basePath, _ := filepath.Abs("./")
	staticPath = path.Join(basePath, staticPath)
	fileServer := http.FileServer(http.Dir(staticPath))
	s.Router.GetInnerRouter().PathPrefix(exposePath).Handler(http.StripPrefix(exposePath, fileServer))
}

func (s *server) GetRouter() Router {
	return s.Router
}

func (s *server) AddHandler(handler Handler) {
	handler.Initialize(s)
	s.Router.AddHandler(handler)
}

func (s *server) ListenAndServe() {
	if err := s.Server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Listen and Server Error: %s", err)
	}
}

func (s *server) Close() {
	go func() {
		if err := s.Server.Shutdown(goContext.TODO()); err != nil {
			panic("panicking" + err.Error())
		}
	}()
	log.Fatal(s.relationalDB.GetDB().Close())
}

// implemented http.Handler interface to our Server so it can be
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		if err := recover(); err != nil {
			handlePanic(err, w, r)
		}
	}()
	s.Router.GetInnerRouter().ServeHTTP(w, r)
	duration := time.Now().Sub(start)
	println(fmt.Sprintf("Request %s %s took %s", r.Method, r.RequestURI, duration.String()))
}

func (s *server) initializeDatabases() {
	if s.relationalDB == nil {
		s.relationalDB = &database.Postgres{
			PSN: s.GetDBConfig().PSN,
		}
		s.relationalDB.InitializeDB()
	}
	if s.cache == nil {
		s.cache = &database.Redis{
			Host:     s.Config.GetRedisConfig().Host,
			Password: s.Config.GetRedisConfig().Password,
		}
		s.cache.InitializeClient()
	}
}

func handlePanic(recoveredPanic interface{}, w http.ResponseWriter, r *http.Request) {
	switch recoveredPanic.(type) {
	case *Error:
		err := recoveredPanic.(*Error)
		logrus.Errorf("Request %s %s returned status %d with stack: \n\n%s\n", r.Method, r.RequestURI, err.StatusCode, debug.Stack())
		w.WriteHeader(err.StatusCode)
		if err.Body != nil {
			if errJson := json.NewEncoder(w).Encode(err.Body); errJson != nil {
				logrus.Errorf("Error %w parsing body %+v when handling an error", errJson, err.Body)
			}
		}
		logError(err)
	default:
		w.WriteHeader(500)
		logrus.Errorf("Request %s %s panicked \"%+v\" with stack: \n\n%s\n", r.Method, r.RequestURI, recoveredPanic, debug.Stack())
	}
}

func logError(error *Error) {
	logrus.Errorf("[STATUS %d] Error %w with body:\n %+v", error.StatusCode, error.Error, error.Body)
}
