package middlewares

import (
	"github.com/gorilla/mux"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"net/http"
)

func UserMiddleware(s server.Server) mux.MiddlewareFunc {
	// userService := authService.New(s.GetRelationalDatabase(), s)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			h.ServeHTTP(w, r)
		})
	}
}
