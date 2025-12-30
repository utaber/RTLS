// dashboard.js - Real-time updates with SSE
console.log('dashboard.js loaded successfully');

// Global variables
let mapInstance = null;
let deviceMarkers = [];
let tempMarker = null;
let eventSource = null;

// Update current time
function updateCurrentTime() {
    const now = new Date();
    const timeString = now.toLocaleTimeString('id-ID', {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
    });
    const dateString = now.toLocaleDateString('id-ID', {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
    
    const currentTimeElement = document.getElementById('currentTime');
    const lastUpdateTimeElement = document.getElementById('lastUpdateTime');
    
    if (currentTimeElement) {
        currentTimeElement.textContent = `${dateString} ${timeString}`;
    }
        
    if (lastUpdateTimeElement) {
        lastUpdateTimeElement.textContent = now.toLocaleTimeString('id-ID', {
            hour: '2-digit', 
            minute: '2-digit'
        });
    }
}

// Update statistics based on devices data
function updateStatistics(devices) {
    const totalDevices = devices.length;
    const activeDevices = devices.filter(d => d.status === "Terdeteksi").length;
    const alertsCount = totalDevices - activeDevices;
    const activePercentage = totalDevices > 0 ? Math.round((activeDevices / totalDevices) * 100) : 0;

    // Update stat cards
    document.querySelectorAll('.card-title').forEach((el, idx) => {
        if (idx === 0) el.textContent = totalDevices;
        if (idx === 1) el.textContent = activeDevices;
        if (idx === 2) el.textContent = activeDevices;
        if (idx === 3) el.textContent = alertsCount;
    });

    // Update small texts
    const smallTexts = document.querySelectorAll('.card-body small.text-muted');
    if (smallTexts[0]) smallTexts[0].textContent = `${activeDevices} active`;
    if (smallTexts[1]) smallTexts[1].textContent = `${activePercentage}% of total`;
    if (smallTexts[3]) smallTexts[3].textContent = alertsCount > 0 ? 'Attention needed' : 'All clear';
}

// Update device table/cards
function updateDeviceTable(devices) {
    // Desktop table
    const tbody = document.querySelector('.d-none.d-md-block tbody');
    if (tbody) {
        if (devices.length === 0) {
            tbody.innerHTML = `
                <tr>
                    <td colspan="4" class="text-center py-5">
                        <div class="text-muted">
                            <i class="bi bi-inbox" style="font-size: 3rem;"></i>
                            <p class="mt-3">No device activity found</p>
                        </div>
                    </td>
                </tr>
            `;
        } else {
            tbody.innerHTML = devices.map(device => `
                <tr>
                    <td>
                        <div class="d-flex align-items-center">
                            <div class="bg-primary bg-opacity-10 p-2 rounded-circle me-3">
                                <i class="bi bi-cpu text-primary"></i>
                            </div>
                            <div>
                                <strong class="d-block">${device.name || 'Unknown'}</strong>
                                <small class="text-muted">${device.deviceID || 'N/A'}</small>
                            </div>
                        </div>
                    </td>
                    <td>
                        <div class="text-monospace">
                            ${device.lat.toFixed(5)}, ${device.lng.toFixed(5)}
                        </div>
                    </td>
                    <td>
                        ${device.status === "Terdeteksi" ? 
                            `<span class="badge bg-success rounded-pill p-2">
                                <i class="bi bi-check-circle me-1"></i>
                                Terdeteksi
                            </span>` :
                            `<span class="badge bg-secondary rounded-pill p-2">
                                <i class="bi bi-x-circle me-1"></i>
                                Tidak Terdeteksi
                            </span>`
                        }
                    </td>
                    <td class="text-center">
                        <button onclick="showOnMap('${device.lat}', '${device.lng}')" 
                                class="btn btn-sm btn-outline-primary">
                            <i class="bi bi-map"></i> View on Map
                        </button>
                    </td>
                </tr>
            `).join('');
        }
    }

    // Mobile cards
    const mobileContainer = document.querySelector('.d-md-none .p-3');
    if (mobileContainer) {
        if (devices.length === 0) {
            mobileContainer.innerHTML = `
                <div class="text-center py-5">
                    <div class="text-muted">
                        <i class="bi bi-inbox" style="font-size: 3rem;"></i>
                        <p class="mt-3 mb-0">No device activity found</p>
                        <small class="d-block">Add devices to see activity here</small>
                        <a href="/devices" class="btn btn-primary mt-3">
                            <i class="bi bi-plus-circle"></i> Add Device
                        </a>
                    </div>
                </div>
            `;
        } else {
            mobileContainer.innerHTML = devices.map(device => `
                <div class="device-card card mb-3 ${device.status === "Terdeteksi" ? 'active' : 'offline'}">
                    <div class="card-body">
                        <div class="d-flex justify-content-between align-items-start mb-2">
                            <div>
                                <h6 class="card-title mb-1">
                                    <i class="bi bi-cpu me-2 text-primary"></i>
                                    ${device.name || 'Unknown'}
                                </h6>
                                <p class="card-text text-muted small mb-2">
                                    <code>${device.deviceID || 'N/A'}</code>
                                </p>
                            </div>
                            <div>
                                ${device.status === "Terdeteksi" ?
                                    `<span class="badge bg-success rounded-pill">
                                        <i class="bi bi-check-circle"></i> Active
                                    </span>` :
                                    `<span class="badge bg-secondary rounded-pill">
                                        <i class="bi bi-x-circle"></i> Offline
                                    </span>`
                                }
                            </div>
                        </div>
                        <div class="mb-3">
                            <p class="mb-1 small text-muted">
                                <i class="bi bi-geo-alt me-1"></i> Location
                            </p>
                            <div class="bg-light p-2 rounded">
                                <div class="text-center">
                                    <strong class="d-block">${device.lat.toFixed(5)}</strong>
                                    <strong class="d-block">${device.lng.toFixed(5)}</strong>
                                </div>
                            </div>
                        </div>
                        <div class="d-flex justify-content-between align-items-center">
                            <div>
                                <small class="text-muted">Device ID: ${device.deviceID}</small>
                            </div>
                            <div>
                                <button onclick="showOnMap('${device.lat}', '${device.lng}')" 
                                        class="btn btn-sm btn-outline-primary">
                                    <i class="bi bi-map"></i>
                                </button>
                                <a href="/devices?device=${device.deviceID}" 
                                   class="btn btn-sm btn-outline-secondary ms-1">
                                    <i class="bi bi-info-circle"></i>
                                </a>
                            </div>
                        </div>
                    </div>
                </div>
            `).join('');
        }
    }
}

// Function to show location on map
function showOnMap(latParam, lngParam) {
    console.log('showOnMap called with:', latParam, lngParam);
    
    const lat = parseFloat(latParam);
    const lng = parseFloat(lngParam);
    
    if (mapInstance && !isNaN(lat) && !isNaN(lng)) {
        mapInstance.setView([lat, lng], 16);
        
        if (tempMarker) {
            mapInstance.removeLayer(tempMarker);
            tempMarker = null;
        }
        
        const blueIcon = L.icon({
            iconUrl: "https://maps.google.com/mapfiles/ms/icons/blue-dot.png",
            iconSize: [40, 40],
            iconAnchor: [20, 40]
        });
        
        tempMarker = L.marker([lat, lng], { icon: blueIcon })
            .addTo(mapInstance)
            .bindPopup(`
                <div class="text-center">
                    <b>Selected Location</b><br>
                    Lat: ${lat.toFixed(5)}<br>
                    Lng: ${lng.toFixed(5)}
                </div>
            `)
            .openPopup();
        
        setTimeout(() => {
            if (tempMarker) {
                mapInstance.removeLayer(tempMarker);
                tempMarker = null;
            }
        }, 8000);
    }
}

// Update map markers
function updateMapMarkers(devices) {
    if (!mapInstance) return;

    // Clear existing markers
    deviceMarkers.forEach(marker => {
        mapInstance.removeLayer(marker);
    });
    deviceMarkers = [];

    const greenIcon = L.icon({
        iconUrl: "https://maps.google.com/mapfiles/ms/icons/green-dot.png",
        iconSize: [32, 32],
        iconAnchor: [16, 32]
    });

    const redIcon = L.icon({
        iconUrl: "https://maps.google.com/mapfiles/ms/icons/red-dot.png",
        iconSize: [32, 32],
        iconAnchor: [16, 32]
    });

    const bounds = [];

    devices.forEach(device => {
        if (!device.lat || !device.lng) return;

        const lat = parseFloat(device.lat);
        const lng = parseFloat(device.lng);
        
        if (isNaN(lat) || isNaN(lng)) return;

        const isActive = device.status === "Terdeteksi";

        const marker = L.marker([lat, lng], { 
            icon: isActive ? greenIcon : redIcon 
        }).addTo(mapInstance);

        marker.bindPopup(`
            <div style="min-width: 200px;">
                <b>${device.name || 'Unknown'}</b><br>
                ID: ${device.deviceID || 'N/A'}<br>
                Status: <span class="badge ${isActive ? 'bg-success' : 'bg-secondary'}">
                    ${device.status || 'Unknown'}
                </span><br>
                Lat: ${lat.toFixed(5)}<br>
                Lng: ${lng.toFixed(5)}
            </div>
        `);

        deviceMarkers.push(marker);
        bounds.push([lat, lng]);
    });

    if (bounds.length > 0) {
        mapInstance.fitBounds(bounds, { padding: [50, 50] });
    }
}

// Initialize map
function initMap() {
    console.log('Initializing map...');
    
    try {
        if (typeof L === 'undefined') {
            console.error('Leaflet library not loaded!');
            return;
        }

        const dataElement = document.getElementById("devices-data");
        if (!dataElement) {
            console.error("devices-data element not found");
            return;
        }

        let devices = [];
        try {
            const dataText = dataElement.textContent.trim();
            const parsedData = JSON.parse(dataText);
            
            if (typeof parsedData === 'string') {
                devices = JSON.parse(parsedData);
            } else {
                devices = parsedData;
            }
            
            if (!Array.isArray(devices)) {
                devices = [];
            }
        } catch (err) {
            console.error("Invalid JSON data", err);
            devices = [];
        }

        const mapElement = document.getElementById("map");
        if (!mapElement) {
            console.error("Map element not found");
            return;
        }
        
        if (!mapInstance) {
            mapInstance = L.map("map").setView([1.1045, 104.0305], 13);
            
            L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
                maxZoom: 19,
                attribution: "© OpenStreetMap"
            }).addTo(mapInstance);
        }

        updateMapMarkers(devices);
        
        console.log("Map initialized successfully");
    } catch (error) {
        console.error("Error initializing map:", error);
    }
}

// Setup SSE connection for real-time updates
function setupRealtimeUpdates() {
    if (typeof EventSource === 'undefined') {
        console.error('SSE not supported by browser');
        return;
    }

    console.log('Setting up real-time updates...');
    
    eventSource = new EventSource('/stream/devices');
    
    eventSource.onopen = function() {
        console.log('SSE connection opened');
        
        // Show connection indicator
        const indicator = document.createElement('div');
        indicator.id = 'sse-indicator';
        indicator.className = 'badge bg-success position-fixed top-0 end-0 m-3';
        indicator.innerHTML = '<i class="bi bi-broadcast"></i> Live';
        document.body.appendChild(indicator);
    };
    
    eventSource.onmessage = function(event) {
        console.log('Received SSE update');
        
        try {
            const devices = JSON.parse(event.data);
            console.log('Devices updated:', devices.length);
            
            // Update all components
            updateStatistics(devices);
            updateDeviceTable(devices);
            updateMapMarkers(devices);
            updateCurrentTime();
            
            // Flash update indicator
            const indicator = document.getElementById('sse-indicator');
            if (indicator) {
                indicator.classList.remove('bg-success');
                indicator.classList.add('bg-warning');
                indicator.innerHTML = '<i class="bi bi-arrow-clockwise"></i> Updating...';
                
                setTimeout(() => {
                    indicator.classList.remove('bg-warning');
                    indicator.classList.add('bg-success');
                    indicator.innerHTML = '<i class="bi bi-broadcast"></i> Live';
                }, 500);
            }
        } catch (err) {
            console.error('Error parsing SSE data:', err);
        }
    };
    
    eventSource.onerror = function(err) {
        console.error('SSE error:', err);
        
        const indicator = document.getElementById('sse-indicator');
        if (indicator) {
            indicator.classList.remove('bg-success');
            indicator.classList.add('bg-danger');
            indicator.innerHTML = '<i class="bi bi-x-circle"></i> Disconnected';
        }
        
        // Attempt to reconnect after 5 seconds
        setTimeout(() => {
            console.log('Attempting to reconnect...');
            eventSource.close();
            setupRealtimeUpdates();
        }, 5000);
    };
}

// Refresh map manually
function refreshMap() {
    console.log('Manual refresh triggered');
    
    if (!mapInstance) {
        initMap();
        return;
    }
    
    mapInstance.invalidateSize();
    
    const dataElement = document.getElementById("devices-data");
    if (dataElement) {
        try {
            const dataText = dataElement.textContent.trim();
            const parsedData = JSON.parse(dataText);
            const devices = typeof parsedData === 'string' ? JSON.parse(parsedData) : parsedData;
            
            if (Array.isArray(devices)) {
                updateMapMarkers(devices);
            }
        } catch (err) {
            console.error("Error parsing device data:", err);
        }
    }
    
    const refreshBtn = document.getElementById('refreshMapBtn');
    if (refreshBtn) {
        const originalHTML = refreshBtn.innerHTML;
        refreshBtn.innerHTML = '<i class="bi bi-check"></i> Refreshed!';
        refreshBtn.classList.remove('btn-outline-primary');
        refreshBtn.classList.add('btn-success');
        
        setTimeout(() => {
            refreshBtn.innerHTML = originalHTML;
            refreshBtn.classList.remove('btn-success');
            refreshBtn.classList.add('btn-outline-primary');
        }, 2000);
    }
}

// Cleanup on page unload
window.addEventListener('beforeunload', function() {
    if (eventSource) {
        eventSource.close();
        console.log('SSE connection closed');
    }
});

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM loaded, initializing dashboard...');
    
    if (typeof L === 'undefined') {
        console.error('ERROR: Leaflet library not loaded!');
        return;
    }
    
    // Initialize map
    initMap();
    
    // Setup real-time updates
    setupRealtimeUpdates();
    
    // Setup refresh button
    const refreshBtn = document.getElementById('refreshMapBtn');
    if (refreshBtn) {
        refreshBtn.addEventListener('click', refreshMap);
    }
    
    // Update time every second
    setInterval(updateCurrentTime, 1000);
    updateCurrentTime();
    
    console.log('Dashboard initialized successfully');
});

// Export functions
window.showOnMap = showOnMap;
window.refreshMap = refreshMap;

console.log('dashboard.js fully loaded');