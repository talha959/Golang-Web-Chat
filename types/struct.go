package types

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

type MyClaims struct {
	Sub string `json:"username.sub"`
	jwt.StandardClaims
}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Recipient string `json:"recipient,omitempty"`
}

var (
	Broadcast       = make(chan Message)
	Connections     = make(map[*websocket.Conn]bool)
	UserConnections = make(map[string]map[*websocket.Conn]bool)
	// Username        = make(map[string]string)
)

// type MyClaims struct {
// 	Sub string `json:"username.sub"`
// 	jwt.StandardClaims
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// type Message struct {
// 	Sender    string `json:"sender"`
// 	Recipient string `json:"recipient,omitempty"`
// 	Message   string `json:"message"`
// }

var (
	broadcast       = make(chan Message)
	connections     = make(map[*websocket.Conn]bool)
	userConnections = make(map[string]map[*websocket.Conn]bool)
)
