package middlewares

import (
	"fmt"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"net/http"
)

func JWTPayload(s server.Server) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: GetJwtSecretFunc(s.GetServerConfig().JWTSecret),
			SigningMethod:       jwt.SigningMethodHS256,
		}).Handler(h)
	}
}

func GetJwtSecretFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}
}
