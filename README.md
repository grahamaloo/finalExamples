# finalExamples
Final project example code

### Lab 9
Here I put in some code for handling login requests sent in via POST
Check out the JS files in static to see how I handle them client side

[Read some examples](https://github.com/gin-gonic/gin#api-examples) to help you out

### otherExample
This example was created by Clinton to illustrate how to do the same thing without using gin, a web framework
This is basically the same, but you have to handle a few extra things yourself.

## General idea:
When you want to retrieve data and display it, use a GET request. When you want to send in data and get a response, use a POST request. You can return HTML or JSON from a POST or GET request.
Other than that use the 4 step example from lab 7 to get it working. You want to query -> add vars -> scan -> display. Then, send back that data to the user

Your code should be something like this:
```
  router.GET("/allTrips", func(c *gin.Context) {
    rows, err := db.Query("YOUR QUERY HERE;")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		// if you are simply inserting data you can stop here. I'd suggest returning a JSON object saying "insert successful" or something along those lines.
		// get all the columns. You can do something with them here if you like, such as adding them to a table header, or adding them to the JSON
		cols, _ := rows.Columns()
		if len(cols) == 0 {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		// This will hold an array of all values
		// makes an array of size 1, storing strings (replace with int or whatever data you want to store)
		output := make([]string, 1)
		
    // The variable(s) here should match your returned columns in the EXACT same order as you give them in your query
		var returnedColumn1 string
		for rows.Next() {
			rows.Scan(&returnedColumn1)
			// VERY important that you store the result back in output
			output = append(output, returnedColumn1)
		}
		//Finally, return your results to the user:
    c.JSON(http.StatusOK, gin.H{"result": output})
  }
```
