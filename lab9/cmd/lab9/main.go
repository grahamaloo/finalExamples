// by setting package as main, Go will compile this as an executable file.
// Any other package turns this into a library
package main

// These are your imports / libraries / frameworks
import (
	// this is Go's built-in sql library
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	// this allows us to run our web server
	"github.com/gin-gonic/gin"
	// this lets us connect to Postgres DB's
	_ "github.com/lib/pq"
)

var (
	// this is the pointer to the database we will be working with
	// this is a "global" variable (sorta kinda, but you can use it as such)
	db *sql.DB
)

func main() {
	rand.Seed(500)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var errd error
	// here we want to open a connection to the database using an environemnt variable.
	// This isn't the best technique, but it is the simplest one for heroku
	db, errd = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if errd != nil {
		log.Fatalf("Error opening database: %q", errd)
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("html/*")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/ping", func(c *gin.Context) {
		ping := db.Ping()
		if ping != nil {
			// our site can't handle http status codes, but I'll still put them in cause why not
			c.JSON(http.StatusOK, gin.H{"error": "true", "message": "db was not created. Contact your TA for assistance"})
		} else {
			c.JSON(http.StatusOK, gin.H{"error": "false", "message": "db created"})
		}
	})

	router.POST("/login1", func(c *gin.Context) {
		// this is meant for SQL injection examples ONLY.
		// Don't copy this for use in an actual environment, even if you do stop SQL injection
		username := c.PostForm("username")
		password := c.PostForm("password")
		if hasIllegalSyntax(username) || hasIllegalSyntax(password) {
			c.JSON(http.StatusOK, gin.H{"result": "failed", "message": "Don't use syntax that isn't allowed"})
			return
		}

		rows, err := db.Query("SELECT usr.name FROM usr WHERE usr.name = '" + username + "' AND usr.password = '" + password + "';")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		cols, _ := rows.Columns()
		if len(cols) == 0 {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		rowCount := 0
		var resultUser string
		for rows.Next() {
			rows.Scan(&resultUser)
			rowCount++
		}
		if rowCount > 1 {
			c.JSON(http.StatusOK, gin.H{"result": "failed", "message": "Too many users returned!"})
			return
		}
		// quick way to check if the user logged in properly
		if rowCount == 0 {
			c.JSON(http.StatusOK, gin.H{"result": "failed", "message": "Wrong password/username!"})
			return
		}

		if resultUser == "admin" {
			c.JSON(http.StatusOK, gin.H{"result": "success", "username": resultUser, "randomCode": rand.Int()})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": "success", "username": resultUser})
		}
	})

	router.POST("/login2", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		if hasIllegalSyntax(username) || hasIllegalSyntax(password) {
			c.JSON(http.StatusOK, gin.H{"result": "failed", "message": "Don't use syntax that isn't allowed"})
			return
		}
		// SQL injection in password only
		rows, err := db.Query("SELECT usr.name FROM usr WHERE usr.name = '" + username + "';")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		cols, _ := rows.Columns()
		if len(cols) == 0 {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		rowCount := 0
		var resultUser string
		for rows.Next() {
			rows.Scan(&resultUser)
			rowCount++
		}
		if rowCount > 1 {
			c.JSON(http.StatusOK, gin.H{"result": "failed", "message": "Too many users returned!"})
			return
		}
		// quick way to check if the user logged in properly
		if rowCount == 0 {
			c.JSON(http.StatusOK, gin.H{"result": "failed", "message": "Wrong password/username!"})
			return
		}

		if resultUser == "admin" {
			c.JSON(http.StatusOK, gin.H{"result": "success", "username": resultUser, "randomCode": rand.Int()})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": "success", "username": resultUser})
		}
	})

	router.POST("/login3", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		if (hasIllegalSyntax(username) || hasIllegalSyntax(password)) && !(strings.Contains(strings.ToUpper(username), "INSERT") || strings.Contains(strings.ToUpper(password), "INSERT")) {
			c.JSON(http.StatusOK, gin.H{"result": "failed", "message": "Don't use syntax that isn't allowed"})
			return
		}
		rows, err := db.Query("SELECT usr.name FROM usr WHERE usr.name = '" + username + "' AND usr.password = '" + password + "';")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		cols, _ := rows.Columns()
		if len(cols) == 0 {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		rowCount := 0
		var resultUser string
		for rows.Next() {
			rows.Scan(&resultUser)
			rowCount++
		}
		if rowCount > 1 {
			c.JSON(http.StatusOK, gin.H{"result": "failed", "message": "Too many users returned!"})
			return
		}
		// quick way to check if the user logged in properly
		if rowCount == 0 {
			c.JSON(http.StatusOK, gin.H{"result": "failed", "message": "Wrong password/username!"})
			return
		}

		if resultUser == "admin" {
			c.JSON(http.StatusOK, gin.H{"result": "success", "username": resultUser, "randomCode": rand.Int()})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": "success", "username": resultUser})
		}
	})

	router.POST("/login4", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		rows, err := db.Query("SELECT usr.name FROM usr WHERE usr.name = $1 AND usr.password = $2;", username, password)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		cols, _ := rows.Columns()
		if len(cols) == 0 {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		rowCount := 0
		var resultUser string
		for rows.Next() {
			rows.Scan(&resultUser)
			rowCount++
		}
		if rowCount > 1 {
			c.JSON(http.StatusOK, gin.H{"result": "failed", "message": "Too many users returned!"})
			return
		}
		// quick way to check if the user logged in properly
		if rowCount == 0 {
			c.JSON(http.StatusOK, gin.H{"result": "failed", "message": "Wrong password/username!"})
			return
		}

		if resultUser != "admin" && resultUser != "Cameron" && resultUser != "normal" && resultUser != "student" {
			c.JSON(http.StatusOK, gin.H{"result": "success", "username": resultUser, "randomCode": rand.Int()})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": "success", "username": resultUser})
		}
	})
	router.POST("/login5", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		rows, err := db.Query("SELECT usr.name FROM usr WHERE usr.name = $1 AND usr.password = $2;", username, password)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		cols, _ := rows.Columns()
		if len(cols) == 0 {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		rowCount := 0
		var resultUser string
		for rows.Next() {
			rows.Scan(&resultUser)
			rowCount++
		}
		if rowCount > 1 {
			c.JSON(http.StatusOK, gin.H{"result": "failed", "message": "Too many users returned!"})
			return
		}
		// quick way to check if the user logged in properly
		if rowCount == 0 {
			c.JSON(http.StatusOK, gin.H{"result": "failed", "message": "Wrong password/username!"})
			return
		}

		if resultUser == "Cameron" {
			c.JSON(http.StatusOK, gin.H{"result": "success", "username": resultUser, "randomCode": rand.Int()})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": "success", "username": resultUser})
		}
	})

	// NO code should go after this line. it won't ever reach that point
	router.Run(":" + port)
}

func hasIllegalSyntax(s string) bool {
	s = strings.ToUpper(s)
	return strings.Contains(s, "INSERT") || strings.Contains(s, "DELETE") || strings.Contains(s, "CREATE") || strings.Contains(s, "DROP") || strings.Contains(s, "UPDATE") || strings.Contains(s, "ALTER")
}
