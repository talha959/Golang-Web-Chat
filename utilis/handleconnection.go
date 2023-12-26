package utilis

import (
	"GIN/types"
	"fmt"
)

func HandleMessage() {
	for {
		message := <-types.Broadcast
		for _, userConns := range types.UserConnections {
			for conn := range userConns {
				err := conn.WriteJSON(message)
				if err != nil {
					fmt.Println(err)
					delete(userConns, conn)
					conn.Close()
				}
			}
		}
	}
}

// var clients = make(map[*websocket.Conn]bool)
// var broadcast = make(chan Message)

// type Message struct {
// 	Username string `json:"username"`
// 	Message  string `json:"message"`
// }

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }
