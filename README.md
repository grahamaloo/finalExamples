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
