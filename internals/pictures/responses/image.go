package responses

type Image struct {
	ImageID    int    `json:"image_id"`
	UserID     int    `json:"user_id"`
	ThumbURL   string `json:"thumb_url"`
	LowResURL  string `json:"lowres_url"`
	HighResURL string `json:"highres_url"`
}
