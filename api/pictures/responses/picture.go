package responses

import (
	"github.com/jopicornell/go-rest-api/api/users/responses"
)

type Picture struct {
	PictureID   int    `json:"picture_id"`
	UserID      int32  `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	NumComments int    `json:"num_comments"`
	NumLikes    int    `json:"num_likes"`
}

type PictureWithImages struct {
	Picture
	Image Image                       `json:"image"`
	User  responses.UserWithoutAvatar `json:"user"`
}
