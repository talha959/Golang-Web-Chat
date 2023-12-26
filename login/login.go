package login

import (
	"GIN/generatetoken"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Loginpage(c *gin.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER")
	passwords := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")
	// table := os.Getenv("DBTABLE")
	fmt.Println(host)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, passwords, dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database!")

	username := c.PostForm("username")
	password := c.PostForm("password")

	var storedUsername, storedPassword string
	err = db.QueryRow("SELECT username, password FROM users WHERE username = $1 AND password = $2", username, password).Scan(&storedUsername, &storedPassword)

	switch {
	case err == sql.ErrNoRows:
		fmt.Println("User not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username not found"})
		return
	case err != nil:
		log.Fatal(err)
		return
	default:
		if password == storedPassword && username == storedUsername {
			fmt.Println("Authentication successful")
			tokenString, err := generatetoken.Token(username, c)
			if err != nil {
				log.Fatal(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating JWT token"})
				return
			}
			fmt.Printf("Generated JWT Token: %s\n", tokenString)
			c.JSON(http.StatusOK, gin.H{
				"JWT":     tokenString,
				"Message": "User login successful",
			})
		} else {
			fmt.Println("Authentication failed: Invalid password")
		}
	}
}
