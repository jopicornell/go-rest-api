package errors

import "errors"

var (
	AuthUserNotMatched = errors.New("user not matched with the given user & password pair")
	UsernameExists     = errors.New("username already exists")
	ForbiddenAction    = errors.New("user is forbidden to do thad action")
)
