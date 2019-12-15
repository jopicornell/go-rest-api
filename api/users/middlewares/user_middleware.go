package middlewares

import (
	"github.com/jopicornell/go-rest-api/api/users/models"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"net/http"
	"strings"
)

type UserMiddleware struct {
	server.Middleware
	Roles []string
}

func (u *UserMiddleware) Handle(res server.Response, req server.Context, next server.HandlerFunc) {
	authHeader := req.GetRequest().Header.Get("Authorization")
	authSplit := strings.Split(authHeader, " ")
	if len(authSplit) == 1 {
		res.Respond(http.StatusUnauthorized)
		return
	}
	authToken := authSplit[1]
	user := new(models.UserWithRoles)
	req.GetServer().GetCache().GetStruct(authToken, user)
	if len(u.Roles) > 0 && !user.HasSomeRole(u.Roles) {
		res.Respond(http.StatusForbidden)
		return
	}
	req.SetUser(user)
	next(res, req)
}
