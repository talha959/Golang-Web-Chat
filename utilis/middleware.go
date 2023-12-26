package utilis

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ExtractUsernameFromToken(c *gin.Context) {

	token := c.GetHeader("Authorization")

	if token != "" {
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Token missing"})
	}
}
