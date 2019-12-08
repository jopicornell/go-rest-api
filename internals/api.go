package internals

import (
	"github.com/jopicornell/go-rest-api/internals/images"
	"github.com/jopicornell/go-rest-api/internals/pictures"
	"github.com/jopicornell/go-rest-api/internals/users"
	"github.com/jopicornell/go-rest-api/pkg/server"
)

func Start() {
	appServer := server.Initialize()
	defer appServer.Close()
	appServer.AddHandler(&users.UserHandler{})
	appServer.AddHandler(&pictures.PictureHandler{})
	appServer.AddHandler(&images.ImageHandler{})
	appServer.AddStatics("/files/images", "./images")
	appServer.AddStatics("", "./static")
	appServer.ListenAndServe()
}
