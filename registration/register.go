package registration

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// const (
// 	host     = os.Getenv("HOST")
// 	port     = ""
// 	user     = ""
// 	password = ""
// 	dbname   = ""
// 	table    = ""
// )

var table = "register"

func Register(c *gin.Context) {
	// host     := os.Getenv("HOST")
	// port:=os.Getenv("PORT")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")

	connStr := "user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection error"})
		return
	}
	tableExists, err := tableExists(db, table)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking if table: " + err.Error()})
		return
	}

	if !tableExists {
		err = createTable(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating table"})
			return
		}
	}
	defer db.Close()
	username := c.PostForm("username")
	passwordValue := c.PostForm("password")
	role := c.PostForm("role")

	_, err = db.Exec("INSERT INTO "+table+" (username, password, role) VALUES ($1, $2, $3);", username, passwordValue, role)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database insertion error", "error status": err})
		// return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "User registration successful",
	})
}
func tableExists(db *sql.DB, tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_name = $1
		)
	`

	var exists bool
	stmt, err := db.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(tableName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func createTable(db *sql.DB) error {
	query := `
		CREATE TABLE ` + table + ` (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL
		)
	`

	_, err := db.Exec(query)
	return err
}
