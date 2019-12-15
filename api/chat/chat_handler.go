package chat

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/jopicornell/go-rest-api/api/users/models"
	"github.com/jopicornell/go-rest-api/pkg/server"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
)

type ChatHandler struct {
	server      server.Server
	upgrader    websocket.Upgrader
	connections []*websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type chatMessage struct {
	MessageType string `json:"type"`
	Content     string `json:"content"`
}

func (a *ChatHandler) Initialize(server server.Server) {
	a.server = server
	a.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	a.connections = []*websocket.Conn{}
}

func (a *ChatHandler) ConfigureRoutes(router server.Router) {
	router.AddRoute("/websocket", a.Subscribe).Methods(http.MethodGet)
}

// Sorry for the mess :D
func (a *ChatHandler) Subscribe(res server.Response, context server.Context) {
	conn, err := upgrader.Upgrade(res.GetWriter(), context.GetRequest(), http.Header{"Sec-WebSocket-Protocol": {"access_token"}})
	if err != nil {
		logrus.Error("Error creating websocket connection", err)
		res.Respond(http.StatusInternalServerError)
	}
	a.connections = append(a.connections, conn)
	authenticated := false
	for {
		messageType, reader, err := conn.NextReader()
		if err != nil {
			logrus.Error("Error processing message", err)
			break
		}
		messageBytes, err := ioutil.ReadAll(reader)

		if err != nil {
			logrus.Errorf("Error reading message from ws: %v", err)
			continue
		}
		var message = new(chatMessage)
		if err := json.Unmarshal(messageBytes, message); err != nil {
			_ = conn.Close()
			logrus.Errorf("Error reading message from ws: %v", err)
			return
		}
		if !authenticated && message.MessageType == "token" {
			authenticated = true
			authToken := message.Content
			user := new(models.UserWithRoles)
			context.GetServer().GetCache().GetStruct(authToken, user)
			if user.UserID == 0 {
				err := conn.Close()
				if err != nil {
					log.Println("Error closing connection", err)
				}
				return
			}
			continue
		}
		if !authenticated {
			_ = conn.Close()
			return
		}
		_ = conn.WriteMessage(messageType, messageBytes)
		for _, otherConn := range a.connections {
			if otherConn != conn {
				_ = otherConn.WriteMessage(messageType, messageBytes)
			}
		}
	}
	_ = conn.Close()

}
