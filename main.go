// main.go

package main

import (
	"GIN/handleconnections"
	"GIN/login"
	"GIN/registration"
	"GIN/types"
	"GIN/utilis"

	// "GIN/"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/register", registration.Register)
	r.POST("/login", login.Loginpage)
	r.Any("/chat", utilis.ExtractUsernameFromToken, utilis.Chat)

	r.Any("/chatting/:ID", func(c *gin.Context) {
		handleconnections.ChatHandler(c, types.Broadcast)
	})

	go utilis.HandleMessage()

	r.Run(":8080")
}
