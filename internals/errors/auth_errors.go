package errors

import "errors"

var AuthUserNotMatched = errors.New("user not matched with the given user & password pair")
