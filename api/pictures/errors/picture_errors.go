package errors

import (
	"errors"
)

var (
	UserAlreadyLikedPicture = errors.New("user has already liked picture")
	PictureNotFound         = errors.New("picture was not found")
	CommentNotFound         = errors.New("comment was not found")
)
