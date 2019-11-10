package handlers

import (
	"github.com/jopicornell/go-rest-api/pkg/config"
	"github.com/jopicornell/go-rest-api/pkg/servertesting"
	"testing"
)

func TestNew(t *testing.T) {
	serverMock := &servertesting.ServerMock{
		Config: config.Config{},
	}
	authHandler := New(serverMock)
	if authHandler == nil {
		t.Errorf("task handler should not be null")
		return
	}
	if authHandler.authService == nil {
		t.Errorf("task handler created without the service")
	}
}

func TestAuthHandler_Login(t *testing.T) {
	t.Run("should throw a bad request if data is invalid", loginBadRequest)
	t.Run("should throw a forbidden when user + password pair is incorrect", loginForbidden)
	t.Run("should return the token when all is correct", loginSuccess)
}

func loginSuccess(t *testing.T) {

}

func loginForbidden(t *testing.T) {

}

func loginBadRequest(t *testing.T) {

}
