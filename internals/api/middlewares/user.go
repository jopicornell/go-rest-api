package middlewares

import (
	"github.com/jopicornell/go-rest-api/pkg/server"
)

type UserMiddleware struct {
	server.Middleware
}

func (u *UserMiddleware) Handle(res server.Response, req server.Request, next server.HandlerFunc) {
	println("handling")
	next(res, req)
}
