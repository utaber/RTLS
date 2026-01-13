// devices.js - Device management functionality

// ========== KONFIGURASI ==========
console.log('devices.js loaded successfully');

// Global variables untuk menyimpan data sementara
let currentEditDeviceId = null;
let currentDeleteDeviceId = null;

// ========== HELPER FUNCTIONS ==========
function showAlertLocal(message, type = 'success') {
    // Coba gunakan showAlert dari main.js jika tersedia
    if (typeof window.showAlert !== 'undefined') {
        window.showAlert(message, type);
    } else {
        // Fallback ke alert browser
        alert(`${type.toUpperCase()}: ${message}`);
    }
}

// ========== DEVICE FUNCTIONS ==========

// Function untuk save device (CREATE)
async function saveDevice() {
    console.log('saveDevice called');
    
    // Cek jika fetchAPI tersedia
    if (typeof window.fetchAPI === 'undefined') {
        showAlertLocal('System not ready. Please wait a moment and try again.', 'danger');
        return;
    }
    
    const name = document.getElementById('add_name').value.trim();

    if (!name) {
        showAlertLocal('Please enter device name', 'danger');
        return;
    }

    try {
        console.log('Creating device with name:', name);
        
        const response = await window.fetchAPI('/api/devices', 'POST', {
            name: name
        });

        console.log('Response status:', response.status);
        
        const result = await response.json();
        console.log('Response result:', result);

        // Handle response
        if (response.ok) {
            const modal = bootstrap.Modal.getInstance(document.getElementById('addDeviceModal'));
            if (modal) modal.hide();

            document.getElementById('addDeviceForm').reset();
            
            // Ambil device_id dari response
            const deviceId = result.device_id || result.data?.device_id || result.id || 'N/A';
            showAlertLocal(`Device "${name}" created successfully with ID: ${deviceId}`, 'success');

            setTimeout(() => {
                window.location.reload();
            }, 1500);
        } else {
            showAlertLocal(result.error || result.detail || result.message || 'Failed to create device', 'danger');
        }
    } catch (error) {
        console.error('Error creating device:', error);
        showAlertLocal('Error creating device. Please try again.', 'danger');
    }
}

// Function untuk UPDATE device
async function updateDevice() {
    console.log('updateDevice called');
    
    // Cek jika fetchAPI tersedia
    if (typeof window.fetchAPI === 'undefined') {
        showAlertLocal('System not ready. Please wait a moment and try again.', 'danger');
        return;
    }
    
    if (!currentEditDeviceId) {
        showAlertLocal('No device selected', 'danger');
        return;
    }

    const name = document.getElementById('edit_name').value.trim();

    if (!name) {
        showAlertLocal('Please enter device name', 'danger');
        return;
    }

    try {
        console.log('Updating device:', currentEditDeviceId, 'with name:', name);
        
        const response = await window.fetchAPI(`/api/devices/${currentEditDeviceId}`, 'PUT', {
            name: name
        });

        console.log('Update response status:', response.status);
        
        const result = await response.json();
        console.log('Update result:', result);

        if (response.ok) {
            const modal = bootstrap.Modal.getInstance(document.getElementById('editDeviceModal'));
            if (modal) modal.hide();

            document.getElementById('editDeviceForm').reset();
            currentEditDeviceId = null;

            showAlertLocal(`Device updated successfully`, 'success');

            setTimeout(() => {
                window.location.reload();
            }, 1500);
        } else {
            showAlertLocal(result.error || result.detail || result.message || 'Failed to update device', 'danger');
        }
    } catch (error) {
        console.error('Error updating device:', error);
        showAlertLocal('Error updating device. Please try again.', 'danger');
    }
}

// Function untuk DELETE device
async function confirmDelete() {
    console.log('confirmDelete called');
    
    // Cek jika fetchAPI tersedia
    if (typeof window.fetchAPI === 'undefined') {
        showAlertLocal('System not ready. Please wait a moment and try again.', 'danger');
        return;
    }
    
    if (!currentDeleteDeviceId) {
        showAlertLocal('No device selected', 'danger');
        return;
    }

    try {
        console.log('Deleting device:', currentDeleteDeviceId);
        
        const response = await window.fetchAPI(`/api/devices/${currentDeleteDeviceId}`, 'DELETE', null);

        console.log('Delete response status:', response.status);
        
        const result = await response.json();
        console.log('Delete result:', result);

        if (response.ok) {
            const modal = bootstrap.Modal.getInstance(document.getElementById('deleteDeviceModal'));
            if (modal) modal.hide();

            currentDeleteDeviceId = null;

            showAlertLocal(`Device deleted successfully`, 'success');

            setTimeout(() => {
                window.location.reload();
            }, 1500);
        } else {
            showAlertLocal(result.error || result.detail || result.message || 'Failed to delete device', 'danger');
        }
    } catch (error) {
        console.error('Error deleting device:', error);
        showAlertLocal('Error deleting device. Please try again.', 'danger');
    }
}

// Function untuk open edit modal
function openEditModal(deviceId, deviceName) {
    console.log('openEditModal called:', deviceId, deviceName);
    
    currentEditDeviceId = deviceId;

    // Set values in modal
    document.getElementById('edit_device_id').value = deviceId;
    document.getElementById('edit_name').value = deviceName;

    // Show modal
    const modal = new bootstrap.Modal(document.getElementById('editDeviceModal'));
    modal.show();
}

// Function untuk open delete modal
function openDeleteModal(deviceId, deviceName) {
    console.log('openDeleteModal called:', deviceId, deviceName);
    
    currentDeleteDeviceId = deviceId;

    // Set values in modal
    document.getElementById('delete_device_id').value = deviceId;
    document.getElementById('delete_name').textContent = deviceName;
    document.getElementById('delete_id').textContent = deviceId;

    // Show modal
    const modal = new bootstrap.Modal(document.getElementById('deleteDeviceModal'));
    modal.show();
}

// ========== INITIALIZATION ==========

// Function untuk initialize setelah semua script siap
function initializeDevices() {
    console.log('Initializing devices functionality...');
    console.log('Auth token cookie:', document.cookie);
    console.log('fetchAPI available:', typeof window.fetchAPI !== 'undefined');
    console.log('Bootstrap available:', typeof bootstrap !== 'undefined');

    // Cek jika semua dependency tersedia
    if (typeof window.fetchAPI === 'undefined') {
        console.error('ERROR: fetchAPI is not available!');
        showAlertLocal('System loading... Please wait a moment.', 'info');
        
        // Coba lagi setelah 1 detik
        setTimeout(initializeDevices, 1000);
        return;
    }

    if (typeof bootstrap === 'undefined') {
        console.error('ERROR: Bootstrap is not available!');
        showAlertLocal('Bootstrap not loaded. Please refresh page.', 'danger');
        return;
    }

    console.log('All dependencies loaded successfully!');

    // Setup event listeners
    setupEventListeners();
    
    console.log('devices.js ready');
}

// Setup event listeners
function setupEventListeners() {
    // Add event listeners untuk form submission dengan Enter key
    const addNameInput = document.getElementById('add_name');
    if (addNameInput) {
        addNameInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                e.preventDefault();
                saveDevice();
            }
        });
    }

    const editNameInput = document.getElementById('edit_name');
    if (editNameInput) {
        editNameInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                e.preventDefault();
                updateDevice();
            }
        });
    }

    // Clear modals when closed
    const addModal = document.getElementById('addDeviceModal');
    if (addModal) {
        addModal.addEventListener('hidden.bs.modal', function() {
            const form = document.getElementById('addDeviceForm');
            if (form) form.reset();
        });
    }

    const editModal = document.getElementById('editDeviceModal');
    if (editModal) {
        editModal.addEventListener('hidden.bs.modal', function() {
            const form = document.getElementById('editDeviceForm');
            if (form) form.reset();
            currentEditDeviceId = null;
        });
    }

    const deleteModal = document.getElementById('deleteDeviceModal');
    if (deleteModal) {
        deleteModal.addEventListener('hidden.bs.modal', function() {
            currentDeleteDeviceId = null;
        });
    }
}

// ========== STARTUP ==========

// Tunggu DOM siap DAN main.js selesai load
function waitForDependencies() {
    console.log('Waiting for dependencies...');
    
    // Cek jika DOM sudah ready
    if (document.readyState === 'loading') {
        console.log('DOM still loading, waiting...');
        document.addEventListener('DOMContentLoaded', waitForDependencies);
        return;
    }
    
    // Cek jika main.js sudah load (fetchAPI tersedia)
    if (typeof window.fetchAPI === 'undefined') {
        console.log('fetchAPI not ready yet, waiting...');
        setTimeout(waitForDependencies, 500);
        return;
    }
    
    // Cek jika Bootstrap sudah load
    if (typeof bootstrap === 'undefined') {
        console.log('Bootstrap not ready yet, waiting...');
        setTimeout(waitForDependencies, 500);
        return;
    }
    
    // Semua dependency siap, initialize
    console.log('All dependencies ready!');
    initializeDevices();
}

// Mulai proses
console.log('Starting devices.js initialization...');
waitForDependencies();

// ========== EXPORT FUNCTIONS ==========

// Export functions untuk global use
window.saveDevice = saveDevice;
window.updateDevice = updateDevice;
window.openEditModal = openEditModal;
window.openDeleteModal = openDeleteModal;
window.confirmDelete = confirmDelete;

console.log('devices.js functions registered for export');