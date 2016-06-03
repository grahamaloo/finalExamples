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
	//var circle = new L.marker([40, -80], {color: 'orange'});
  //circle.addTo(address);
  address.addTo(map);
  var address1 = "United States";

  var defaultStyle = {
    color: "#2262CC",
    weight: 2,
    opacity: 0.6,
    fillOpacity: 0.1,
    fillColor: "#2262CC"
  };
  var hoverStyle = {
    color: '#2262CC', 
    weight: 3,
    opacity: 0.6,
    fillOpacity: 0.65,
    fillColor: '#2262CC'
  };
  
  var onEachFeature = function(feature, layer) {
    layer.setStyle(defaultStyle);

    (function(layer, properties) {
      layer.on("mouseover", function (e) {
        layer.setStyle(hoverStyle);
      });

      layer.on("mouseout", function (e) {
        layer.setStyle(defaultStyle);
      });
    })(layer, feature.properties);
  };

  var geoJSONlayer = L.geoJson({
    onEachFeature : onEachFeature
  }).addTo(map);

  $.get("/districts", function(data) {
    console.log(data);
    for (var i = 0; i < data.districts.length; i++) {
      data.districts[i] = JSON.parse(data.districts[i]);
      geoJSONlayer.addData(data.districts[i]);
    }
    console.log(data);
  })
	//console.log(address1.latlng);
	//L.esri.Geocoding.geocode().address('380 New York St').city('Redlands').region('California').postal(92373).run(function(err, results, response){
  //		var circle2 = new L.marker([results.results[0].latlng.lat,results.results[0].latlng.lng]);
  //		circle2.addTo(map);
	//});

  var arcgisOnline = L.esri.Geocoding.arcgisOnlineProvider();

  var searchControl = L.esri.Geocoding.geosearch({
    providers: [
      arcgisOnline,
      new L.esri.Geocoding.MapServiceProvider({
        label: 'States and Counties',
        url: 'https://sampleserver6.arcgisonline.com/arcgis/rest/services/Census/MapServer',
        layers: [2, 3],
        searchFields: ['NAME', 'STATE_NAME']
      })
    ]
  }).addTo(map);

  searchControl.on("results", function(data) {
     map.setView([data.results[0].latlng.lat,data.results[0].latlng.lng], 7);
  });
	
	$.get("/myquery", function(data){
		console.log(data);
  });

  $.get("/addresses", function(data){
    //console.log(data);
    //console.log(data.Addresses.length);
    //console.log(data.Addresses[0].LineOne);
    for (var i = 0; i < data.Addresses.length; i++) {
      //console.log("inside loop");
      L.esri.Geocoding.geocode().address(data.Addresses[i].LineOne + " " + data.Addresses[i].LineTwo).city(data.Addresses[i].City).region(data.Addresses[i].State).run(function(err, results, response){
        var circle2 = new L.marker([results.results[0].latlng.lat,results.results[0].latlng.lng]);
        console.log(results.results[0].latlng.lat);
        circle2.addTo(map);
      })
    }
  });

  var popup = L.popup();

  function onMapClick(e) {
    popup
      .setLatLng(e.latlng)
      .setContent("You clicked the map at " + e.latlng.toString())
      .openOn(map);
}

map.on('click', onMapClick);

  $("#submit1").click(function(){
      if($('input[name="prv"]:checked').val() == 1) {
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
})



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
      $.post("/donationOldPerson", {email: $("#email-old").val(), amount: $("#amount-old").val(), payment: $("#payment-id-old").val()}).done(function(data) {
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
      $.post("/donationNewPersonCard", {email: $("#email-new").val(), amount: $("#amount-new").val(), payment: $("#payment-id-new").val(),
      								f_name: $("#f_name").val(), l_name: $("#l_name").val(), phone: $("#phone").val(), addr_line_1: $("#addr-line-1").val(),
      								addr_line_2: $("#addr-line-2").val(), city: $("#city").val(), state_code: $("#state-code").val(),
      								cardNumber: $("#card-num-new").val(), cardExp: $("#exp-new").val()}).done(function(data) {
        if(data.result == "failed") {
          $("#result-old").text("" + data.message);
        } else {
          console.log(data);
          $("#result-old").text("Success! " + data.message);
        }
      })
    }

    function donationNewPerson() {
      $.post("/donationNewPerson", {email: $("#email-new").val(), amount: $("#amount-new").val(), payment: $("#payment-id-new").val(),
      								f_name: $("#f_name").val(), l_name: $("#l_name").val(), phone: $("#phone").val(), addr_line_1: $("#addr-line-1").val(),
      								addr_line_2: $("#addr-line-2").val(), city: $("#city").val(), state_code: $("#state-code").val()}).done(function(data) {
        if(data.result == "failed") {
          $("#result-old").text("" + data.message);
        } else if (data.result == "succeeded"){
          $("#result-old").text("Success! " + data.message);
        }
      })
    }

