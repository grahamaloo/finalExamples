$(function(){

	var map = L.map('map').setView([38.152, -95.625], 4);
     //Create a tile layer variable using the appropriate url
    L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token={accessToken}', {
    attribution: 'Map data &copy; <a href="http://openstreetmap.org">OpenStreetMap</a> contributors, <a href="http://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery © <a href="http://mapbox.com">Mapbox</a>',
    maxZoom: 18,   
    id: 'mapbox.dark',
    accessToken: 'pk.eyJ1Ijoia2V2aW5ka2UiLCJhIjoiY2lmdTF1MDV3MWlmOHQ1bHl3bmgyYXUwcCJ9.8Givw8o8IVmz9n6Gckshkg'
	}).addTo(map);
	
	var address = new L.LayerGroup();
	var circle = new L.marker([40, -80], {color: 'orange'});
    circle.addTo(address);
    address.addTo(map);
    var address1 = "United States";
	//console.log(address1.latlng);
	L.esri.Geocoding.geocode().address('380 New York St').city('Redlands').region('California').postal(92373).run(function(err, results, response){
  		var circle2 = new L.marker([results.results[0].latlng.lat,results.results[0].latlng.lng]);
  		console.log(results.results[0].latlng.lat);
  		circle2.addTo(map);
	});
	
	$.get("/myquery", function(data){
		console.log(data);
  });
  $.get("/addresses", function(data){
    console.log(data);
  });
	
/*
	$.get("/addresses", function(data){
        $("#firstQuery").append(data);
    }, "html")
  */  
    /*
    	geocoder = new google.maps.Geocoder();
  		geocoder.geocode( {address:address}, function(results, status) {
		alert(results[0].geometry.lat());
     	 var marker = new L.circleMarker([results[0].geometry.lat(),results[0].geometry.lat()], {color: 'orange'});
  });
  	marker.addTo(map);
  	*/
})
	//var thumbsUp = element(by.css('span.glyphicon-thumbs-up'));
	//var thumbsDown = element(by.css('span.glyphicon-thumbs-down'));

/*
    $.get("/ping", function(data){
        if(data.error == "true"){
            $("#results").prepend("<div class='alert alert-danger'><strong>Error!</strong> "+ data.message +"</div>");
        }
    }, "json")

    $("#submit1").click(function(){
      login(1);
    });
    $("#submit2").click(function(){
      login(2);
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


