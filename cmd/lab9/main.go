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
	"time"
	//"strconv"
	//"strings"

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

	/*router.POST("/ping", func(c *gin.Context) {
		ping := db.Ping()
		if ping != nil {
			// our site can't handle http status codes, but I'll still put them in cause why not
			c.JSON(http.StatusOK, gin.H{"error": "true", "message": "db was not created. Contact your TA for assistance"})
		} else {
			c.JSON(http.StatusOK, gin.H{"Ned", "Caetlyn", "Rob", "Ygritte", "Osha", "Hodor"});
		}
	})*/

	type NameList struct {
		Names []string
	}

	type Address struct {
		LineOne string
		LineTwo sql.NullString
		City string
		State string
	}

	type AddressList struct {
		Addresses []Address
	}
	
	router.GET("/myquery", func(c *gin.Context) {
		names := NameList{[]string{"Ned", "Caetlyn", "Rob", "Ygritte", "Osha", "Hodor"}}
		//c.JSON(http.StatusOK, gin.H{"names":[]interface{}{"Ned", "Caetlyn", "Rob", "Ygritte", "Osha", "Hodor"},});
		c.JSON(http.StatusOK, names)
	})

	router.GET("/addresses", func(c *gin.Context) {
		rows, err := db.Query("SELECT first_line, second_line, city, state_code FROM address;")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		cols, _ := rows.Columns()
		if len(cols) == 0 {
			c.AbortWithStatus(http.StatusNoContent)
		}
		var tempAddresses []Address
		var first_line string
		var second_line sql.NullString
		var city string
		var state_code string
		for rows.Next() {
			rows.Scan(&first_line, &second_line, &city, &state_code)
			add := Address{first_line, second_line, city, state_code}
			tempAddresses = append(tempAddresses, add)
		}

		c.JSON(http.StatusOK, AddressList{tempAddresses})
		})

	router.POST("/donationOldPerson", func(c *gin.Context) {
		email := c.PostForm("email")
		amount := c.PostForm("amount")
		paymentId := c.PostForm("payment") // assume for now payment is passing the id. this is not normal functionality

		var personId int64
		err := db.QueryRow("SELECT person.person_id FROM person WHERE person.email = $1;", email).Scan(&personId)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"result":"failed", "message":"no user with that email"})
			return
		}
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		current_time := time.Now().Local()
		_, err = db.Exec("SELECT add_donation($1, $2, $3, $4);" , amount, current_time.Format("2006-01-02"), personId, paymentId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	})
	
	router.POST("/donationOldPersonCard", func(c *gin.Context) {
		email := c.PostForm("email")
		amount := c.PostForm("amount")
		paymentId := c.PostForm("payment") // assume for now payment is passing the id. this is not normal functionality
		//card_num := c.PostForm("cardNumber")
		//card_exp := c.PostForm("cardExp")

		var personId int64
		err := db.QueryRow("SELECT person.person_id FROM person WHERE person.email = $1;", email).Scan(&personId)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"result":"failed", "message":"no user with that email"})
			return
		}
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		current_time := time.Now().Local()
		_, err = db.Exec("SELECT add_donation($1, $2, $3, $4);" , amount, current_time.Format("2006-01-02"), personId, paymentId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	})
	
	router.POST("/donationNewPerson", func(c *gin.Context) {
		email := c.PostForm("email")
		amount := c.PostForm("amount")
		paymentId := c.PostForm("payment") // assume for now payment is passing the id. this is not normal functionality
		f_name := c.PostForm("f_name")
		l_name := c.PostForm("l_name")
		phone := c.PostForm("phone")
		addr_line_1 := c.PostForm("addr_line_1")
		addr_line_2 := c.PostForm("addr_line_2")
		city := c.PostForm("city")
		state_code := c.PostForm("state_code")


		//_, err := db.Exec("SELECT insert_person($1, $2, $3, $4, $5, $6, $7, $8);", f_name, l_name, phone, email, addr_line_1, addr_line_2, city, state_code)
		_, err := db.Exec("SELECT insert_person('Graham', 'Kelly', '2068904595', 'gereae', '214 dsaf a', 'dddd', 'Seattle', 'WA');")
		if err != nil {
			//c.JSON(http.StatusOK, gin.H{"result":"failed", "message":"person insert did not succeed at part 1"})
			//c.JSON
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		
		var personId int64
		err = db.QueryRow("SELECT person.person_id FROM person WHERE person.email = $1;", email).Scan(&personId)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"result":"failed", "message":"person insert did not succeed"})
			return
		}
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		
		current_time := time.Now().Local()
		_, err = db.Exec("SELECT add_donation($1, $2, $3, $4)", amount, current_time.Format("2006-01-02"), personId, paymentId)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"result":"failed", "message":"donation insert did not succeed"})
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"result":"succeeded", "message":"donation and user successfully added (well no errors at least)"})
			return
		}
	})
	router.POST("/donationNewPersonCard", func(c *gin.Context) {
		email := c.PostForm("email")
		amount := c.PostForm("amount")
		paymentId := c.PostForm("payment") // assume for now payment is passing the id. this is not normal functionality
		f_name := c.PostForm("f_name")
		l_name := c.PostForm("l_name")
		phone := c.PostForm("phone")
		addr_line_1 := c.PostForm("addr_line_1")
		addr_line_2 := c.PostForm("addr_line_2")
		city := c.PostForm("city")
		state_code := c.PostForm("state_code")
		//card_num := c.PostForm("cardNumber")
		//card_exp := c.PostForm("cardExp")

		_, err := db.Exec("SELECT insert_person($1, $2, $3, $4, $5, $6, $7, $8);", f_name, l_name, phone, email, addr_line_1, addr_line_2, city, state_code)
		
		if err != nil {
			//c.JSON
			//c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		
		var personId int64
		err = db.QueryRow("SELECT person.person_id FROM person WHERE person.email = $1;", email).Scan(&personId)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"result":"failed", "message":"person insert did not succeed"})
			return
		}
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		
		current_time := time.Now().Local()
		_, err = db.Exec("SELECT add_donation($1, $2, $3, $4)", amount, current_time.Format("2006-01-02"), personId, paymentId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"result":"succeeded", "message":"donation and user successfully added (well no errors at least)"})
			return
		}
	})
/*
	router.POST("/submit1", func(c *gin.Context) {
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

	router.POST("/submit2", func(c *gin.Context) {
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
	
	router.GET("/addresses", func(c *gin.Context) {
			//var a := [2]string{4506 NE 17th ave, Seattle Washington, 8007 NE 179th Place Seattle Washington}
			db.Query("SELECT a.first_line,a.second_line, address FROM address AS a NATURAL JOIN Car WHERE Car.brand = 'Honda' AND Car.model = 'Civic'")
	})
	*/

	

	// NO code should go after this line. it won't ever reach that point
	router.Run(":" + port)
}


