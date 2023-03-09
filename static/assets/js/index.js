function initMap() {
    const map = new google.maps.Map(document.getElementById("map"), { zoom: 10, center: { lat: -33.9, lng: 151.2 }, }); setMarkers(map);
}    

function initMap() {  
    const myLatLng = { lat: -25.363, lng: 131.044 };  
    const map = new google.maps.Map(document.getElementById("map"), 
    {    zoom: 4,    center: myLatLng,  });  
    new google.maps.Marker({    position: myLatLng,    map,    title: "Hello World!", 
 });
}
window.initMap = initMap;