package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type router struct {
	router *mux.Router
	server Server
}

type Router interface {
	http.Handler
	AddGroup(path string) Router
	AddRoute(path string, handler HandlerFunc) *mux.Route
	Use(...Middleware) Router
	GetInnerRouter() *mux.Router
	AddHandler(handler Handler)
}

func NewRouter(server Server) Router {
	return &router{router: mux.NewRouter(), server: server}
}

func HandleHTTP(handler HandlerFunc, server Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(NewResponse(w), NewContext(r, server))
	}
}

func (r *router) AddHandler(handler Handler) {
	handler.ConfigureRoutes(r)
}

func (r *router) AddGroup(path string) Router {
	return &router{router: r.router.PathPrefix(path).Subrouter(), server: r.server}
}

func (r *router) AddRoute(path string, handler HandlerFunc) *mux.Route {
	return r.router.HandleFunc(path, HandleHTTP(handler, r.server))
}

func (r *router) Use(middlewares ...Middleware) Router {
	convertedMiddlewares := make([]mux.MiddlewareFunc, len(middlewares))
	for index, middleware := range middlewares {
		convertedMiddlewares[index] = func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				middleware.Handle(NewResponse(w), NewContext(req, r.server), func(res Response, req Context) {
					handler.ServeHTTP(res.GetWriter(), req.GetRequest())
				})
			})
		}
	}
	r.router.Use(convertedMiddlewares...)
	return r
}

func (r *router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(res, req)
}

func (r *router) GetInnerRouter() *mux.Router {
	return r.router
}
