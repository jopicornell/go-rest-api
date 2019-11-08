package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/pkg/config"
	"github.com/jopicornell/go-rest-api/pkg/database"
	"log"
	"net/http"
	"path"
	"path/filepath"
)

func Initialize() Server {
	server := &server{}
	server.Config.Bootstrap()
	server.initializeRelationalDatabase()
	server.Router = mux.NewRouter().StrictSlash(true)
	server.ApiRouter = server.Router.PathPrefix(server.GetServerConfig().ApiUrl).Subrouter()
	return server
}

type Server interface {
	Close()
	GetRelationalDatabase() *sqlx.DB
	AddApiRoute(path string, handler HandlerFunc) *mux.Route
	AddRoute(string, HandlerFunc) *mux.Route
	AddStatics(string, string)
	ListenAndServe()
}

type server struct {
	config.Config
	relationalDB *database.MySQL
	Router       *mux.Router
	ApiRouter    *mux.Router
}

func (s *server) Close() {
	log.Fatal(s.relationalDB.GetDB().Close())
}

func (s *server) GetRelationalDatabase() *sqlx.DB {
	return s.relationalDB.GetDB()
}

func (s *server) AddApiRoute(path string, handler HandlerFunc) *mux.Route {
	return s.ApiRouter.HandleFunc(path, HandleHTTP(handler))
}

func (s *server) AddRoute(path string, handler HandlerFunc) *mux.Route {
	return s.Router.HandleFunc(path, HandleHTTP(handler))
}

func (s *server) AddStatics(exposePath string, staticPath string) {
	basePath, _ := filepath.Abs("./")
	staticPath = path.Join(basePath, staticPath)
	fileServer := http.FileServer(http.Dir(staticPath))
	s.Router.PathPrefix(exposePath).Handler(http.StripPrefix(exposePath, fileServer))
}

func (s server) ListenAndServe() {
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", s.GetServerConfig().Port), s.Router))
}

func (s *server) initializeRelationalDatabase() *sqlx.DB {
	if s.relationalDB == nil {
		s.relationalDB = &database.MySQL{
			PSN: s.GetDBConfig().PSN,
		}
	}
	return s.relationalDB.GetDB()
}
