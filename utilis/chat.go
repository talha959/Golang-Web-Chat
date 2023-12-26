package utilis

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type MyClaims struct {
	Sub string `json:"username.sub"`
	jwt.StandardClaims
}

func Chat(c *gin.Context) {
	token := c.GetHeader("Authorization")
	splitToken := strings.Split(token, ".")
	tokenlength := len(splitToken)

	if len(splitToken) != 3 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Message": "Invalid token format",
			"tokens":  tokenlength,
		})
		return
	}
	// parsedToken, err := jwt.ParseWithClaims(splitToken[1], func(t *jwt.Token) (interface{}, error) {
	// 	return []byte("AllYourBase"), nil
	// })

	// if err != nil || !parsedToken.Valid {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"Message": "Failed to parse or invalid token",
	// 		"token":   parsedToken,
	// 	})
	// 	return
	// }
	// claims, ok := parsedToken.Claims.(jwt.MapClaims)
	// if !ok {
	// 	fmt.Errorf("invalid token claims")
	// // }
	tokens, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Provide the key used to sign the token
		var secretKey = []byte("AllYourBase")
		ss, err := token.SignedString(secretKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error"})
			return "", err
		}

		return ss, nil
	})
	if err != nil {
		fmt.Println("error")
	}
	claims, ok := tokens.Claims.(*MyClaims)
	if !ok {
		fmt.Println("claims error")

	}
	fmt.Println(claims)
	c.JSON(http.StatusOK, gin.H{
		"Message":  "Welcome to the chat",
		"username": claims.Subject,
	})
}
