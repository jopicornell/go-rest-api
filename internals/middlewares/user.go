package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"net/http"
)

func UserMiddleware(s server.Server) mux.MiddlewareFunc {
	// userService := authService.New(s.GetRelationalDatabase(), s)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Context().Value("user").(*jwt.Token)
			claims, _ := token.Claims.(jwt.MapClaims)
			fmt.Println(claims["sub"])
			h.ServeHTTP(w, r)
		})
	}
}
