# Final Project

This is the code for the final database project for INFO 340 at the Univeristy of Washington. Our database was designed to support a grassroots political campaign. The code here queries the database to display information about the location of events on a map, and can write new donations into the database. The database itself is hosted on heroku and a live version of this application can be found at: group-2-final.herokuapp.com.

### Some Notes:

1. It would be very easy to change the query in the go file to accept different state parameters. i.e. a quick modification of that file and of the JS would allow you to permit the user to specify for which states they'd like to see the districts displayed.

2. Many of the queries rely on stored procedures that cannot be modified, so some of the flexibility here is limited (e.g. it is annoying to select the GeoJSON data from the DB without the function called in the go file, but this function only allows you to specify a state to select from, further granularity in terms of what is returned is not supported or easy to implement).

3. There is far more funtionality supported by our database, but these are the only features we had time to get working (the design and implementation on the DB end was, by far, the most important part of this project).
