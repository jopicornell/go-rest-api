package server

import (
	"github.com/gorilla/mux"
	"net/http"
	"path"
	"path/filepath"
)

type router struct {
	router *mux.Router
}

type Router interface {
	AddGroup(path string) Router
	AddRoute(path string, handler HandlerFunc) *mux.Route
	Use(...mux.MiddlewareFunc)
	GetInnerRouter() *mux.Router
	AddStatics(exposePath string, staticPath string)
}

func NewRouter() Router {
	return &router{router: mux.NewRouter()}
}

func (r *router) AddGroup(path string) Router {
	return &router{router: r.router.PathPrefix(path).Subrouter()}
}

func (r *router) AddRoute(path string, handler HandlerFunc) *mux.Route {
	return r.router.HandleFunc(path, HandleHTTP(handler))
}

func (r *router) Use(middlewares ...mux.MiddlewareFunc) {
	r.router.Use(middlewares...)
}

func (r *router) AddStatics(exposePath string, staticPath string) {
	basePath, _ := filepath.Abs("./")
	staticPath = path.Join(basePath, staticPath)
	fileServer := http.FileServer(http.Dir(staticPath))
	r.router.PathPrefix(exposePath).Handler(http.StripPrefix(exposePath, fileServer))
}

func (r *router) GetInnerRouter() *mux.Router {
	return r.router
}
