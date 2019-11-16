package server

import (
	goContext "context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/pkg/config"
	"github.com/jopicornell/go-rest-api/pkg/database"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"runtime/debug"
	"time"
)

func Initialize() Server {
	server := &server{}
	server.Config.Bootstrap()
	server.initializeRelationalDatabase()
	server.Router = NewRouter()
	server.Server = http.Server{
		Addr:    ":" + server.GetServerConfig().Port,
		Handler: server,
	}
	return server
}

type Server interface {
	http.Handler
	Close()
	GetRelationalDatabase() *sqlx.DB
	GetServerConfig() *config.Server
	GetRouter() Router
	AddStatics(string, string)
	ListenAndServe()
}

type server struct {
	http.Server
	config.Config
	relationalDB *database.MySQL
	Router       Router
}

func (s *server) Close() {
	go func() {
		if err := s.Server.Shutdown(goContext.TODO()); err != nil {
			panic("panicking" + err.Error())
		}
	}()
	log.Fatal(s.relationalDB.GetDB().Close())
}

func (s *server) GetRelationalDatabase() *sqlx.DB {
	return s.relationalDB.GetDB()
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

func (s *server) ListenAndServe() {
	if err := s.Server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Listen and Server Error: %s", err)
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logPanic(err, r)
		}
	}()
	s.Router.GetInnerRouter().ServeHTTP(w, r)
	duration := time.Now().Sub(start)
	println(fmt.Sprintf("Request %s %s took %s", r.Method, r.RequestURI, duration.String()))
}

func (s *server) initializeRelationalDatabase() *sqlx.DB {
	if s.relationalDB == nil {
		s.relationalDB = &database.MySQL{
			PSN: s.GetDBConfig().PSN,
		}
	}
	return s.relationalDB.GetDB()
}

func logPanic(recoveredPanic interface{}, r *http.Request) {
	log.Printf("Request %s %s panicked \"%+v\" with stack: \n\n%s\n", r.Method, r.RequestURI, recoveredPanic, debug.Stack())
}
