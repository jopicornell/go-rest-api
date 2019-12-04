package server

type Middleware interface {
	Handle(Response, Context, HandlerFunc)
}
