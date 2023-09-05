package upgrader

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//允许所有来源的WebSocket连接
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
