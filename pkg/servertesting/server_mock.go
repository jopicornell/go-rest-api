package servertesting

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/pkg/config"
	"github.com/jopicornell/go-rest-api/pkg/database"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"time"
)

func Initialize(cfg *config.Config) *ServerMock {
	serverMock := new(ServerMock)
	serverMock.Config = *cfg
	serverMock.Router = server.NewRouter()
	serverMock.Server = http.Server{
		Addr:    ":" + serverMock.GetServerConfig().Port,
		Handler: serverMock,
	}
	return serverMock
}

type ServerMock struct {
	http.Server
	config.Config
	relationalDB *database.MySQL
	Router       server.Router
}

func (s *ServerMock) Close() {
	log.Fatal(s.relationalDB.GetDB().Close())
}

func (s *ServerMock) GetServerConfig() *config.Server {
	return &config.Server{
		ApiUrl:      "",
		StaticsPath: "",
		Port:        "",
		JWTSecret:   "secret",
	}
}

func (s *ServerMock) GetRelationalDatabase() *sqlx.DB {
	db, _, _ := sqlmock.New()
	return sqlx.NewDb(db, "mock")
}

func (s *ServerMock) GetRouter() server.Router {
	return s.Router
}

func (s *ServerMock) AddHandler(handler server.Handler) {
	s.Router.AddHandler(handler)
}

func (s *ServerMock) AddStatics(exposePath string, staticPath string) {
	basePath, _ := filepath.Abs("./")
	staticPath = path.Join(basePath, staticPath)
	fileServer := http.FileServer(http.Dir(staticPath))
	s.Router.GetInnerRouter().PathPrefix(exposePath).Handler(http.StripPrefix(exposePath, fileServer))
}

func (s *ServerMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	s.Router.GetInnerRouter().ServeHTTP(w, r)
	duration := time.Now().Sub(start)
	println(fmt.Sprintf("Request %s %s took %s", r.Method, r.RequestURI, duration.String()))
}

func (s *ServerMock) ListenAndServe() {
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", s.GetServerConfig().Port), s))
}
