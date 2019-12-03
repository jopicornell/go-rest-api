package server

type HandlerFunc func(Response, Request)

type Handler interface {
	Initialize(Server)
	ConfigureRoutes(Router)
}
