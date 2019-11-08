package servertesting

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/pkg/config"
	"github.com/jopicornell/go-rest-api/pkg/database"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"log"
	"net/http"
	"path"
	"path/filepath"
)

func Initialize(cfg *config.Config) *ServerMock {
	serverMock := new(ServerMock)
	serverMock.Config = *cfg
	serverMock.Router = mux.NewRouter().StrictSlash(true)
	serverMock.ApiRouter = serverMock.Router.PathPrefix(serverMock.GetServerConfig().ApiUrl).Subrouter()
	return serverMock
}

type ServerMock struct {
	config.Config
	relationalDB *database.MySQL
	Router       *mux.Router
	ApiRouter    *mux.Router
}

func (s *ServerMock) Close() {
	log.Fatal(s.relationalDB.GetDB().Close())
}

func (s *ServerMock) GetServerConfig() *config.Server {
	return &config.Server{
		ApiUrl:      "",
		StaticsPath: "",
		Port:        "",
	}
}

func (s *ServerMock) GetRelationalDatabase() *sqlx.DB {
	db, _, _ := sqlmock.New()
	return sqlx.NewDb(db, "mock")
}

func (s *ServerMock) AddApiRoute(path string, handler server.HandlerFunc) *mux.Route {
	return s.ApiRouter.HandleFunc(path, server.HandleHTTP(handler))
}

func (s *ServerMock) AddRoute(path string, handler server.HandlerFunc) *mux.Route {
	return s.Router.HandleFunc(path, server.HandleHTTP(handler))
}

func (s *ServerMock) AddStatics(exposePath string, staticPath string) {
	basePath, _ := filepath.Abs("./")
	staticPath = path.Join(basePath, staticPath)
	fileServer := http.FileServer(http.Dir(staticPath))
	s.Router.PathPrefix(exposePath).Handler(http.StripPrefix(exposePath, fileServer))
}

func (s ServerMock) ListenAndServe() {
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", s.GetServerConfig().Port), s.Router))
}
