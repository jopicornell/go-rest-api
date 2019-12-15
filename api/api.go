package api

import (
	"github.com/jopicornell/go-rest-api/api/chat"
	"github.com/jopicornell/go-rest-api/api/images"
	"github.com/jopicornell/go-rest-api/api/pictures"
	"github.com/jopicornell/go-rest-api/api/users"
	"github.com/jopicornell/go-rest-api/pkg/server"
)

func Start() {
	appServer := server.Initialize()
	defer appServer.Close()
	appServer.AddHandler(&users.UserHandler{})
	appServer.AddHandler(&pictures.PictureHandler{})
	appServer.AddHandler(&images.ImageHandler{})
	appServer.AddHandler(&chat.ChatHandler{})
	appServer.AddStatics("/files/images", "./images")
	appServer.AddStatics("", "./static")
	appServer.ListenAndServe()
}
