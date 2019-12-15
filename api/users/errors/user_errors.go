package errors

import "errors"

var (
	AuthUserNotMatched = errors.New("user not matched with the given user & password pair")
	UsernameExists     = errors.New("username already exists")
	UserNotFound       = errors.New("user not found")
	ForbiddenAction    = errors.New("user is forbidden to do thad action")
)
