// utils.go
package utilis

import (
	"GIN/types"
	"encoding/json"
	"fmt"
)

type PrivateMessage struct {
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

// Private handles private messages
func Private(messageJSON string) error {
	var privateMessage PrivateMessage
	err := json.Unmarshal([]byte(messageJSON), &privateMessage)
	if err != nil {
		return err
	}
	// for r := range types.Username {
	// 	fmt.Println(r, "nnn")
	// }
	recipientUserID := privateMessage.Recipient
	message := privateMessage.Message
	fmt.Println(len(types.UserConnections))
	// fmt.Println(string(types.UserConnections))
	for r := range types.UserConnections {
		fmt.Println(r, "ttt")
	}
	connections, ok := types.UserConnections[recipientUserID]
	if !ok {
		return fmt.Errorf("Recipient not found: %s", recipientUserID)
	}
	// claims, ok := types.MyClaims.Subject

	for conn := range connections {
		err := conn.WriteJSON(map[string]string{"sender": "private", "message": message})
		if err != nil {
			fmt.Println(err)
			delete(types.UserConnections[recipientUserID], conn)
			conn.Close()
		}
	}

	return nil
}
