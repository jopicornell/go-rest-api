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

func Initialize() *Server {
	server := &Server{}
	server.Config.Bootstrap()
	server.initializeRelationalDatabase()
	server.Router = mux.NewRouter().StrictSlash(true)
	server.ApiRouter = server.Router.PathPrefix(server.GetServerConfig().ApiUrl).Subrouter()
	return server
}

type Server struct {
	config.Config
	relationalDB *database.MySQL
	Router       *mux.Router
	ApiRouter    *mux.Router
}

func (s *Server) Close() {
	log.Fatal(s.relationalDB.GetDB().Close())
}

func (s *Server) GetRelationalDatabase() *sqlx.DB {
	return s.relationalDB.GetDB()
}

func (s *Server) AddApiRoute(path string, handler HandlerFunc) *mux.Route {
	return s.ApiRouter.HandleFunc(path, HandleJSONResponse(handler))
}

func (s *Server) AddRoute(path string, handler HandlerFunc) *mux.Route {
	return s.Router.HandleFunc(path, HandleJSONResponse(handler))
}

func (s *Server) AddStatics(exposePath string, staticPath string) {
	basePath, _ := filepath.Abs("./")
	staticPath = path.Join(basePath, staticPath)
	fileServer := http.FileServer(http.Dir(staticPath))
	s.Router.PathPrefix(exposePath).Handler(http.StripPrefix(exposePath, fileServer))
}

func (s *Server) AddRouteWithSerializerFunc(path string, handler HandlerFunc, serializer HandlerSerializer) {
	s.Router.HandleFunc(path, serializer(handler))
}

func (s Server) ListenAndServe() {
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%s", s.GetServerConfig().Port), s.Router))
}

func (s *Server) initializeRelationalDatabase() *sqlx.DB {
	if s.relationalDB == nil {
		s.relationalDB = &database.MySQL{
			PSN: s.GetDBConfig().PSN,
		}
	}
	return s.relationalDB.GetDB()
}
