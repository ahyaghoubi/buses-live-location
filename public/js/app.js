const map = L.map('map');

map.setView([34.794188, 48.486974], 16);

L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '© OpenStreetMap'
}).addTo(map);

const busIcon = L.icon({
    iconUrl: './assets/png/bus1.png', // Replace with the path or URL to your bus icon image
    iconSize: [32, 37], // Adjust the size as needed
    iconAnchor: [16, 37], // The point of the icon which corresponds to the marker's location
    popupAnchor: [0, -37] // The point from which popups will "open", relative to the iconAnchor
})

let markers = []

function mapChanger(busId, lat, lng) {
    markers = markers.filter(markerObj => {
        if (markerObj.busId === busId) {
            map.removeLayer(markerObj.marker);
            return false;
        }
        return true;
    })

    const marker = L.marker([lat, lng], { icon: busIcon }).addTo(map)

    marker.bindPopup(`<strong>Bus ID:</strong> ${busId}`)

    markers.push({
        busId,
        marker
    })
}

const host = window.location.host
const clientWS = new WebSocket("ws://" + host + "/clientws")

clientWS.onopen = function () {
    console.log("Connected to the WebSocket server.");
}

clientWS.onmessage = function (event) {
    const dataArray = JSON.parse(event.data)
    dataArray.forEach((data) => {
        mapChanger(data.bus_id, data.latitude, data.longitude)
    })
}

clientWS.onerror = function (error) {
    console.error("WebSocket error:", error);
}

fetch("http://" + host + "/bus")
    .then((response) => {
        if (!response.ok) {
            throw new Error("Network response was not ok: " + response.status)
        }
        return response.json()
    })
    .then((dataArray) => {
        if (dataArray != null) {
            dataArray.forEach((data) => {
                mapChanger(data.bus_id, data.latitude, data.longitude)
            })
        }
    })
    .catch((error) => {
        console.error("There was a problem with the fetch operation:", error)
    })
