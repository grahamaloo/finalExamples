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
	
	// for testing purposes only
	router.GET("/myquery", func(c *gin.Context) {
		names := NameList{[]string{"Ned", "Caetlyn", "Rob", "Ygritte", "Osha", "Hodor"}}
		c.JSON(http.StatusOK, names)
	})

	// returns all addresses associated with events in DB via GET
	router.GET("/addresses", func(c *gin.Context) {
		//rows, err := db.Query("SELECT first_line, second_line, city, state_code FROM address;")
		rows, err := db.Query("SELECT a.first_line, a.second_line, a.city, a.state_code FROM address AS a NATURAL JOIN event AS e;")
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

	router.GET("/districts", func(c *gin.Context) {
		
		})

	// inserts a new donation for an old person
	router.POST("/donationOldPerson", func(c *gin.Context) {
		email := c.PostForm("email")
		amount := c.PostForm("amount")

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

		var paymentId int64
		err = db.QueryRow("WITH A AS (INSERT INTO payment_method VALUES (DEFAULT) RETURNING payment_method.payment_method_id) INSERT INTO check_payment(payment_method_id) VALUES ((SELECT payment_method_id FROM A)) RETURNING check_payment.payment_method_id;").Scan(&paymentId)
		_, err = db.Exec("SELECT add_donation($1, $2, $3, $4);" , amount, current_time.Format("2006-01-02"), personId, paymentId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"result":"succeeded", "message":"donation and user successfully added (well no errors at least)"})
			return
		}
	})
	
	// insert a new donation for an old person using a card. it will look for the card in the DB before inserting a new card number.
	router.POST("/donationOldPersonCard", func(c *gin.Context) {
		email := c.PostForm("email")
		amount := c.PostForm("amount")
		card_num := c.PostForm("cardNumber") 
		card_exp := c.PostForm("cardExp")

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

		var paymentId int64
		err = db.QueryRow("SELECT credit_card_payment.card_number FROM credit_card_payment WHERE credit_card_payment.card_number = $1;", card_num).Scan(&card_num)
		if err == sql.ErrNoRows {
			err = db.QueryRow("WITH A AS (INSERT INTO payment_method VALUES (DEFAULT) RETURNING payment_method.payment_method_id) INSERT INTO credit_card_payment(payment_method_id, card_number, exp) VALUES((SELECT payment_method_id FROM A), $1,$2) RETURNING credit_card_payment.payment_method_id;", card_num, card_exp).Scan(&paymentId)
		} else if err != nil {
			//c.AbortWithError(http.StatusInternalServerError, err)
			c.JSON(http.StatusOK, gin.H{"result":"failed", "message":"card num look up failed"})
			return
		} else {
			err = db.QueryRow("SELECT ccp.payment_method_id FROM credit_card_payment AS ccp WHERE ccp.card_number = $1", card_num).Scan(&paymentId)
			if err != nil {
			//c.AbortWithError(http.StatusInternalServerError, err)
			c.JSON(http.StatusOK, gin.H{"result":"failed", "message":"select payment id failed", "card_num": card_num, "err": err})
			return
			}
		}

		current_time := time.Now().Local()
		_, err = db.Exec("SELECT add_donation($1, $2, $3, $4);" , amount, current_time.Format("2006-01-02"), personId, paymentId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"result":"succeeded", "message":"donation and user successfully added (well no errors at least)"})
			return
		}
	})
	
	// insert a new donation for a new person
	router.POST("/donationNewPerson", func(c *gin.Context) {
		email := c.PostForm("email")
		amount := c.PostForm("amount") 
		f_name := c.PostForm("f_name")
		l_name := c.PostForm("l_name")
		phone := c.PostForm("phone")
		addr_line_1 := c.PostForm("addr_line_1")
		addr_line_2 := c.PostForm("addr_line_2")
		city := c.PostForm("city")
		state_code := c.PostForm("state_code")


		_, err := db.Exec("SELECT insert_person($1, $2, $3, $4, $5, $6, $7, $8)", f_name, l_name, phone, email, addr_line_1, addr_line_2, city, state_code)
		//_, err := db.Exec("SELECT insert_person('Graham', 'Kelly', '2068904595', 'grahamtk@uw.edu', '214 dsaf a', 'dddd', 'Seattle', 'WA');")
		if err != nil {
			//c.JSON(http.StatusOK, gin.H{"result":"failed", "message":"person insert did not succeed at part 1"})
			//c.JSON
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		
		var personId int64
		err = db.QueryRow("SELECT person.person_id FROM person WHERE person.email = $1", email).Scan(&personId)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"result":"failed", "message":"person insert did not succeed"})
			return
		}
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		
		current_time := time.Now().Local()
		
		var paymentId int64
		err = db.QueryRow("WITH A AS (INSERT INTO payment_method VALUES (DEFAULT) RETURNING payment_method.payment_method_id) INSERT INTO check_payment(payment_method_id) VALUES ((SELECT payment_method_id FROM A)) RETURNING check_payment.payment_method_id;").Scan(&paymentId)

		_, err = db.Exec("SELECT add_donation($1, $2, $3, $4)", amount, current_time.Format("2006-01-02"), personId, paymentId)
		
		if err != nil {
			//c.JSON(http.StatusOK, gin.H{"result":"failed", "message":"donation insert did not succeed"})
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"result":"succeeded", "message":"donation and user successfully added (well no errors at least)"})
			return
		}
	})

	// insert new donation for a new person. this will still look for the card they used in the db before creating a new entry.
	router.POST("/donationNewPersonCard", func(c *gin.Context) {
		email := c.PostForm("email")
		amount := c.PostForm("amount")
		f_name := c.PostForm("f_name")
		l_name := c.PostForm("l_name")
		phone := c.PostForm("phone")
		addr_line_1 := c.PostForm("addr_line_1")
		addr_line_2 := c.PostForm("addr_line_2")
		city := c.PostForm("city")
		state_code := c.PostForm("state_code")
		card_num := c.PostForm("cardNumber") // assume for now payment is passing the id. this is not normal functionality
		card_exp := c.PostForm("cardExp")

		_, err := db.Exec("SELECT insert_person($1, $2, $3, $4, $5, $6, $7, $8);", f_name, l_name, phone, email, addr_line_1, addr_line_2, city, state_code)
		
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"result":"failed", "message":"failed at insert person"})
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
			//c.AbortWithError(http.StatusInternalServerError, err)
			c.JSON(http.StatusOK, gin.H{"result":"failed", "message":"failed at select person id"})
			return
		}
		
		var paymentId int64
		err = db.QueryRow("SELECT credit_card_payment.card_number FROM credit_card_payment WHERE credit_card_payment.card_number = $1;", card_num).Scan(&card_num)
		if err == sql.ErrNoRows {
			err = db.QueryRow("WITH A AS (INSERT INTO payment_method VALUES (DEFAULT) RETURNING payment_method.payment_method_id) INSERT INTO credit_card_payment(payment_method_id, card_number, exp) VALUES((SELECT payment_method_id FROM A), $1,$2) RETURNING credit_card_payment.payment_method_id;", card_num, card_exp).Scan(&paymentId)
		} else if err != nil {
			//c.AbortWithError(http.StatusInternalServerError, err)
			c.JSON(http.StatusOK, gin.H{"result":"failed", "message":"card num look up failed"})
			return
		} else {
			err = db.QueryRow("SELECT ccp.payment_method_id FROM credit_card_payment AS ccp WHERE ccp.card_number = $1", card_num).Scan(&paymentId)
			if err != nil {
			//c.AbortWithError(http.StatusInternalServerError, err)
			c.JSON(http.StatusOK, gin.H{"result":"failed", "message":"select payment id failed", "card_num": card_num, "err": err})
			return
			}
		}

		current_time := time.Now().Local()

		_, err = db.Exec("SELECT add_donation($1, $2, $3, $4);", amount, current_time.Format("2006-01-02"), personId, paymentId)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"result":"succeeded", "message":"donation and user successfully added (well no errors at least)"})
			return
		}
	})

	// NO code should go after this line. it won't ever reach that point
	router.Run(":" + port)
}


