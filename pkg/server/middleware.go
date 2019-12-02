package server

type Middleware interface {
	Handle(Response, Request, HandlerFunc)
}
