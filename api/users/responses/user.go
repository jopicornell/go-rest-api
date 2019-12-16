package responses

import imageResponses "github.com/jopicornell/go-rest-api/api/images/responses"

type User struct {
	UserID   int32                 `json:"user_id"`
	Username string                `json:"username"`
	FullName string                `json:"fullname"`
	Image    *imageResponses.Image `json:"avatar"`
}

type UserWithoutAvatar struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
	FullName string `json:"fullname"`
}
