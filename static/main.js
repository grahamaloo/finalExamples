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
  		circle2.addTo(map);
	});

    var arcgisOnline = L.esri.Geocoding.arcgisOnlineProvider();

    var searchControl = L.esri.Geocoding.geosearch({
    providers: [
      arcgisOnline,
      new L.esri.Geocoding.MapServiceProvider({
        label: 'States and Counties',
        url: 'https://sampleserver6.arcgisonline.com/arcgis/rest/services/Census/MapServer',
        layers: [4, 4],
        searchFields: ['NAME', 'STATE_NAME']
      })
    ]
  }).addTo(map);
	
	$.get("/myquery", function(data){
		console.log(data);
  });

  $.get("/addresses", function(data){
    console.log(data);
    console.log(data.Addresses.length);
    console.log(data.Addresses[0].LineOne);
    for (var i = 0; i < data.Addresses.length; i++) {
      console.log("inside loop");
      L.esri.Geocoding.geocode().address(data.Addresses[i].LineOne + " " + data.Addresses[i].LineTwo).city(data.Addresses[i].City).region(data.Addresses[i].State).run(function(err, results, response){
        var circle2 = new L.marker([results.results[0].latlng.lat,results.results[0].latlng.lng]);
        console.log(results.results[0].latlng.lat);
        circle2.addTo(map);
      })
    }
  });

  $("#submit1").click(function(){
      if($('input[name="prv"]:checked').val() == 1) {
        alert("submitted");
         donationOldPersonCard();
      } else {
         donationOldPerson();
      }
    });
    
    $("#submit2").click(function(){
      if($('input[name="new"]:checked').val() == 1) {
         donationNewPersonCard();
      } else {
         donationNewPerson();
      }
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

    function donationOldPersonCard() {
      $.post("/donationOldPersonCard", {email: $("#email-old").val(), amount: $("#amount-old").val(), payment: $("#payment-id-old").val()
                                            , cardNumber: $("#card-num-old").val(), cardExp: $("#exp-old").val()}).done(function(data) {
        if(data.result == "failed") {
          console.log(data);
          $("#result-old").text("" + data.message);
        } else {
          console.log(data);
          $("#result-old").text("Success! " + data.message);
        }
      })
    }
        function donationOldPerson() {
      $.post("/donationOldPerson", {email: $("#email-old").val(), amount: $("#amount_old").val(), payment: $("#payment-id-old").val()}).done(function(data) {
        if(data.result == "failed") {
          console.log(data);
          $("#result-old").text("" + data.message);
        } else {
          console.log(data);
          $("#result-old").text("Success! " + data.message);
        }
      })
    }
        function donationNewPersonCard() {
      $.post("/donationNewPersonCard", {email: $("#mail-new").val(), amount: $("#amount-new").val(), payment: $("#payment-id-new").val(),
      								f_name: $("#f_name").val(), l_name: $("#l_name").val(), phone: $("#phone").val(), addr_line_1: $("#addr-line-1").val(),
      								addr_line_2: $("#addr-line-2").val(), city: $("#city").val(), state_code: $("#state-code").val(),
      								cardNumber: $("#card-num-new").val(), cardExp: $("#exp-new").val()}).done(function(data) {
        if(data.result == "failed") {
          alert("failed");
          $("#result-old").text("" + data.message);
        } else {
          alert("success...");
          console.log(data);
          $("#result-old").text("Success! " + data.message);
        }
      })
    }
    function donationNewPerson() {
    console.log({email: $("#email-new").val(), amount: $("#amount-new").val(), payment: $("#payment-id-new").val(),
      								f_name: $("#f_name").val(), l_name: $("#l_name").val(), phone: $("#phone").val(), addr_line_1: $("#addr-line-1").val(),
      								addr_line_2: $("#addr-line-2").val(), city: $("#city").val(), state_code: $("#state-code").val()});
      $.post("/donationNewPerson", {email: $("#email-new").val(), amount: $("#amount-new").val(), payment: $("#payment-id-new").val(),
      								f_name: $("#f_name").val(), l_name: $("#l_name").val(), phone: $("#phone").val(), addr_line_1: $("#addr-line-1").val(),
      								addr_line_2: $("#addr-line-2").val(), city: $("#city").val(), state_code: $("#state-code").val()}).done(function(data) {
        if(data.result == "failed") {
          alert("failed");
          console.log(data);
          $("#result-old").text("" + data.message);
        } else if (data.result == "succeeded"){
          alert("success...");
          console.log(data);
          $("#result-old").text("Success! " + data.message);
        }
      })
    }

