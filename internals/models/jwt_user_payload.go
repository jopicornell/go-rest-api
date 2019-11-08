package models

import "github.com/gbrlsnchs/jwt/v3"

type JwtUserPayload struct {
	jwt.Payload
	ID int `json:"id"`
}
