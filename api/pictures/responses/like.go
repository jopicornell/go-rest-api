package responses

type Like struct {
	PictureID int32 `json:"picture_id"`
	UserID    int32 `json:"user_id"`
}
