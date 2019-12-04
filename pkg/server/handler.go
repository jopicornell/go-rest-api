package server

type HandlerFunc func(Response, Context)

type Handler interface {
	Initialize(Server)
	ConfigureRoutes(Router)
}
