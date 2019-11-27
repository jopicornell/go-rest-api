package handlers

import (
	"testing"
)

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
