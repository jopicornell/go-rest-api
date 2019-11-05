package errors

import "errors"

var InternalServerError = errors.New("Internal Server Error")
var NotFound = errors.New("Data not found")
var BadRequest = errors.New("Bad Request")
