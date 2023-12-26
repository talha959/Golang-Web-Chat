package handleconnections

import (
	"GIN/types"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/websocket"
)

func ChatHandler(c *gin.Context, broadcast chan types.Message) {
	userID := c.Param("ID")

	token := c.GetHeader("Authorization")
	tokens, err := jwt.ParseWithClaims(token, &types.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Use a secure way to generate the secret key
		return []byte("AllYourBase"), nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := tokens.Claims.(*types.MyClaims)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid claims"})
		return
	}

	ws, err := types.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer ws.Close()

	types.Connections[ws] = true

	// Add the connection to the userConnections map
	if types.UserConnections[userID] == nil {
		types.UserConnections[userID] = make(map[*websocket.Conn]bool)
	}
	types.UserConnections[userID][ws] = true

	// Write claims inside the loop
	err = ws.WriteJSON(gin.H{"claims": claims})
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		mt, message, err := ws.ReadMessage()
		fmt.Println(mt)
		if err != nil {
			fmt.Println(err)
			break
		}

		if string(message) == "ping" {
			message = []byte("pong")
		}

		// Check if the message is intended for private messaging
		var msg types.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Println(err)
			continue
		}

		if msg.Recipient != "" {
			recipientID := msg.Recipient
			if types.UserConnections[recipientID] != nil {
				for conn := range types.UserConnections[recipientID] {
					err := conn.WriteJSON(types.Message{Sender: claims.Subject, Message: msg.Message})
					if err != nil {
						fmt.Println(err)
						delete(types.UserConnections[recipientID], conn)
						conn.Close()
					}
				}
			}
		} else {
			// Broadcast message only to the connections associated with the user ID
			for conn := range types.UserConnections[userID] {
				err := conn.WriteJSON(types.Message{Sender: claims.Subject, Message: msg.Message})
				if err != nil {
					fmt.Println(err)
					delete(types.UserConnections[userID], conn)
					conn.Close()
				}
			}
		}
	}
}
