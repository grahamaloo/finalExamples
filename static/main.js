$(function(){


	var map = L.map('map').setView([38.152, -95.625], 4);
     //Create a tile layer variable using the appropriate url
    L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token={accessToken}', {
    attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery © <a href="http://mapbox.com">Mapbox</a>',
    maxZoom: 18,   
    id: 'mapbox.dark',
    accessToken: 'pk.eyJ1Ijoia2V2aW5ka2UiLCJhIjoiY2lmdTF1MDV3MWlmOHQ1bHl3bmgyYXUwcCJ9.8Givw8o8IVmz9n6Gckshkg'
	}).addTo(map);
/*
    $.get("/ping", function(data){
        if(data.error == "true"){
            $("#results").prepend("<div class='alert alert-danger'><strong>Error!</strong> "+ data.message +"</div>");
        }
    }, "json")

    $("#login1").click(function(){
      login(1);
    });
    $("#login2").click(function(){
      login(2);
    });
    $("#login3").click(function(){
      login(3);
    });
    $("#login4").click(function(){
      login(4);
    });
    $("#login5").click(function(){
      login(5);
    });

    function login(index){
      $.post("/login"+index, {username: $("#username"+index).val(), password: $("#password"+index).val()})
        .done(function(data){
          if(data.result == "failed"){
            console.log(data)
            $("#result"+index).text("Failed to login! " + data.message);
          } else {
            console.log(data)
            $("#result"+index).text("Logged in as: " + data.username + (data.randomCode ? " (CODE: " + data.randomCode + ")" : ""));
          }
        });
    }
    */

})
