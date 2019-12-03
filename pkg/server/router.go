package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type router struct {
	router *mux.Router
}

type Router interface {
	http.Handler
	AddGroup(path string) Router
	AddRoute(path string, handler HandlerFunc) *mux.Route
	Use(...Middleware)
	GetInnerRouter() *mux.Router
	AddHandler(handler Handler)
}

func NewRouter() Router {
	return &router{router: mux.NewRouter()}
}

func HandleHTTP(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(NewResponse(w), NewRequest(r))
	}
}

func (r *router) AddHandler(handler Handler) {
	handler.ConfigureRoutes(r)
}

func (r *router) AddGroup(path string) Router {
	return &router{router: r.router.PathPrefix(path).Subrouter()}
}

func (r *router) AddRoute(path string, handler HandlerFunc) *mux.Route {
	return r.router.HandleFunc(path, HandleHTTP(handler))
}

func (r *router) Use(middlewares ...Middleware) {
	convertedMiddlewares := make([]mux.MiddlewareFunc, len(middlewares))
	for index, middleware := range middlewares {
		convertedMiddlewares[index] = func(handler http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				middleware.Handle(NewResponse(w), NewRequest(r), func(res Response, req Request) {
					handler.ServeHTTP(res.GetWriter(), req.GetRequest())
				})
			})
		}
	}
	r.router.Use(convertedMiddlewares...)
}

func (r *router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(res, req)
}

func (r *router) GetInnerRouter() *mux.Router {
	return r.router
}
