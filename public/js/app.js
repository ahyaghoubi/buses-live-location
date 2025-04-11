const map = L.map('map');
map.setView([34.794188, 48.486974], 16);

L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '© OpenStreetMap'
}).addTo(map);

let buses = {};
let activeBusId = null;

const busIcon = L.icon({
    iconUrl: './assets/png/bus1.png', // Replace with the path or URL to your bus icon image
    iconSize: [32, 37], // Adjust the size as needed
    iconAnchor: [16, 37], // The point of the icon which corresponds to the marker's location
    popupAnchor: [0, -37] // The point from which popups will "open", relative to the iconAnchor
});

function formatCoordinate(coord) {
    return typeof coord === 'number' && !isNaN(coord) ? coord.toFixed(4) : 'نامشخص';
}

function updateBusMarker(busId, lat, lng) {
    if (isNaN(lat) || isNaN(lng)) return;

    if (buses[busId]?.marker) {
        map.removeLayer(buses[busId].marker);
    }

    const marker = L.marker([lat, lng], {
        icon: busIcon
    }).addTo(map);

    marker.bindPopup(`
        <div class="w-40">
          <h3 class="font-bold text-blue-600">اتوبوس ${busId}</h3>
          <p class="text-xs text-gray-500 mt-1">آخرین بروزرسانی: ${new Date().toLocaleTimeString('fa-IR')}</p>
        </div>
      `);

    buses[busId] = { marker, lat, lng };

    updateBusListItem(busId, lat, lng);
}

function updateBusListItem(busId, lat, lng) {
    let listItem = document.getElementById(`bus-${busId}`);

    if (!listItem) {
        listItem = document.createElement('div');
        listItem.id = `bus-${busId}`;
        listItem.className = 'bus-item p-3 flex items-center cursor-pointer hover:bg-gray-50';
        listItem.innerHTML = `
          <div class="w-8 h-8 rounded-full flex items-center justify-center mr-3 bg-blue-100 text-blue-800">
            <i class="fas fa-bus text-sm"></i>
          </div>
          <div class="flex-1">
            <h3 class="font-medium">اتوبوس ${busId}</h3>
            <p class="text-xs text-gray-500">${formatCoordinate(lat)}, ${formatCoordinate(lng)}</p>
          </div>
          <i class="fas fa-chevron-left text-gray-400 ml-2"></i>
        `;

        listItem.addEventListener('click', () => {
            document.querySelectorAll('.bus-item').forEach(item => item.classList.remove('active'));
            listItem.classList.add('active');
            map.setView([lat, lng], 16);
            buses[busId].marker.openPopup();
            activeBusId = busId;
        });

        document.getElementById('busList').prepend(listItem);
    } else {
        const coordsEl = listItem.querySelector('p.text-gray-500');
        coordsEl.textContent = `${formatCoordinate(lat)}, ${formatCoordinate(lng)}`;
    }

    document.getElementById('lastUpdated').textContent = new Date().toLocaleTimeString('fa-IR');
}

const host = window.location.host;
let clientWS;

function connectWebSocket() {
    clientWS = new WebSocket("ws://" + host + "/clientws");

    clientWS.onmessage = function (event) {
        try {
            const dataArray = JSON.parse(event.data);
            document.getElementById('busList').innerHTML = '';
            if (!Array.isArray(dataArray)) throw new Error("Received data is not an array");

            dataArray.forEach((data) => {
                if (!data.bus_id || data.latitude === undefined || data.longitude === undefined) return;
                const lat = parseFloat(data.latitude);
                const lng = parseFloat(data.longitude);
                updateBusMarker(data.bus_id, lat, lng);
            });
        } catch (error) {
            console.error("WebSocket message error:", error);
        }
    };

    clientWS.onclose = function () {
        setTimeout(connectWebSocket, 5000);
    };
}

function fetchBusData() {
    fetch("http://" + host + "/bus")
        .then(res => res.json())
        .then(dataArray => {
            if (dataArray !== null && dataArray.length > 0) {
                document.getElementById('busList').innerHTML = '';
                if (!Array.isArray(dataArray)) throw new Error("Invalid bus data");

                dataArray.forEach((data) => {
                    if (!data.bus_id || data.latitude === undefined || data.longitude === undefined) return;
                    const lat = parseFloat(data.latitude);
                    const lng = parseFloat(data.longitude);
                    updateBusMarker(data.bus_id, lat, lng);
                });
            }
        })
        .catch(error => {
            console.error("Fetch error:", error);
            document.getElementById('busList').innerHTML = `
            <div class="text-center py-10 text-red-500">
              <i class="fas fa-exclamation-triangle text-3xl mb-2"></i>
              <p>خطا در بارگذاری اطلاعات اتوبوس‌ها</p>
              <p class="text-xs mt-2">${error.message}</p>
            </div>
          `;
        });
}

connectWebSocket();
fetchBusData();

document.querySelector('input[type="text"]').addEventListener('input', function () {
    const searchTerm = this.value.toLowerCase();
    document.querySelectorAll('.bus-item').forEach(item => {
        const busId = item.id.replace('bus-', '');
        item.style.display = busId.includes(searchTerm) || item.textContent.toLowerCase().includes(searchTerm)
            ? 'flex' : 'none';
    });
});

const sidebar = document.getElementById('sidebar');
const sidebarToggle = document.getElementById('sidebarToggle');
const closeSidebar = document.getElementById('closeSidebar');

function toggleSidebar() {
    sidebar.classList.toggle('collapsed');
    sidebarToggle.classList.toggle('expanded');
    sidebarToggle.classList.toggle('collapsed');
    sidebarToggle.innerHTML = sidebar.classList.contains('collapsed')
        ? '<i class="fas fa-chevron-left"></i>'
        : '<i class="fas fa-chevron-right"></i>';
}

sidebarToggle.addEventListener('click', toggleSidebar);
closeSidebar.addEventListener('click', toggleSidebar);