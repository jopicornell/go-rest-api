package models

type Token struct {
	Token  string   `json:"token"`
	UserID int      `json:"user_id"`
	Roles  []string `json:"roles"`
}
