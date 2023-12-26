package utilis

import (
	"GIN/types"
	"fmt"
)

func HandleMessages() {
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
