package server

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/jopicornell/go-rest-api/pkg/config"
	"github.com/jopicornell/go-rest-api/pkg/database"
	"log"
	"net/http"
)

func Initialize() *Server {
	server := &Server{}
	server.Config.Bootstrap()
	server.initializeRelationalDatabase()
	server.Router = mux.NewRouter().StrictSlash(true)
	return server
}

type Server struct {
	config.Config
	relationalDB *database.MySQL
	Router       *mux.Router
}

func (s *Server) Close() {
	log.Fatal(s.relationalDB.GetDB().Close())
}

func (s *Server) GetRelationalDatabase() *sqlx.DB {
	return s.relationalDB.GetDB()
}

func (s *Server) AddRoute(path string, handler HandlerFunc) {
	s.Router.HandleFunc(path, HandleJSONResponse(handler))
}

func (s *Server) AddRouteWithSerializerFunc(path string, handler HandlerFunc, serializer HandlerSerializer) {
	s.Router.HandleFunc(path, serializer(handler))
}

func (s Server) ListenAndServe() {
	log.Panic(http.ListenAndServe(":8080", s.Router))
}

func (s *Server) initializeRelationalDatabase() *sqlx.DB {
	if s.relationalDB == nil {
		s.relationalDB = &database.MySQL{
			PSN: s.GetDBConfig().PSN,
		}
	}
	return s.relationalDB.GetDB()
}
