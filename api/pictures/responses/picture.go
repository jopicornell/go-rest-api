package responses

type Picture struct {
	PictureID   int    `json:"picture_id"`
	UserID      int    `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	NumComments int    `json:"num_comments"`
	NumLikes    int    `json:"num_likes"`
}

type PictureWithImages struct {
	Picture
	Image Image `json:"image"`
}
