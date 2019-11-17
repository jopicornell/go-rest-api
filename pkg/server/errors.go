package server

type Error struct {
	StatusCode int
	Body       interface{}
	Error      error
}
