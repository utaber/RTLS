// devices.js - Device management functionality
console.log('devices.js loaded successfully');

// Global variables untuk menyimpan data sementara
let currentEditDeviceId = null;
let currentDeleteDeviceId = null;

// Function untuk show alert
function showAlert(message, type = 'success') {
    const alertContainer = document.getElementById('alertContainer');
    if (!alertContainer) return;

    const alertHTML = `
        <div class="alert alert-${type} alert-dismissible fade show" role="alert">
            <i class="bi bi-${type === 'success' ? 'check-circle' : 'exclamation-triangle'} me-2"></i>
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        </div>
    `;

    alertContainer.innerHTML = alertHTML;

    // Auto hide after 5 seconds
    setTimeout(() => {
        const alert = alertContainer.querySelector('.alert');
        if (alert) {
            alert.classList.remove('show');
            setTimeout(() => {
                alertContainer.innerHTML = '';
            }, 150);
        }
    }, 5000);
}

// Function untuk save device (Create)
async function saveDevice() {
    console.log('saveDevice called');
    
    const name = document.getElementById('add_name').value.trim();

    if (!name) {
        showAlert('Please enter device name', 'danger');
        return;
    }

    try {
        const response = await fetch('/api/devices', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: name
            })
        });

        const result = await response.json();

        if (response.ok && result.success) {
            // Close modal
            const modal = bootstrap.Modal.getInstance(document.getElementById('addDeviceModal'));
            if (modal) modal.hide();

            // Clear form
            document.getElementById('addDeviceForm').reset();

            // Show success alert
            showAlert(`Device "${name}" created successfully with ID: ${result.device_id}`, 'success');

            // Reload page after short delay
            setTimeout(() => {
                window.location.reload();
            }, 1500);
        } else {
            showAlert(result.error || 'Failed to create device', 'danger');
        }
    } catch (error) {
        console.error('Error creating device:', error);
        showAlert('Error creating device. Please try again.', 'danger');
    }
}

// Function untuk open edit modal
async function openEditModal(deviceId, deviceName) {
    console.log('openEditModal called:', deviceId, deviceName);
    
    currentEditDeviceId = deviceId;

    // Set values in modal
    document.getElementById('edit_device_id').value = deviceId;
    document.getElementById('edit_name').value = deviceName;

    // Show modal
    const modal = new bootstrap.Modal(document.getElementById('editDeviceModal'));
    modal.show();
}

// Function untuk update device
async function updateDevice() {
    console.log('updateDevice called');
    
    if (!currentEditDeviceId) {
        showAlert('No device selected', 'danger');
        return;
    }

    const name = document.getElementById('edit_name').value.trim();

    if (!name) {
        showAlert('Please enter device name', 'danger');
        return;
    }

    try {
        const response = await fetch(`/api/devices/${currentEditDeviceId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: name
            })
        });

        const result = await response.json();

        if (response.ok && result.success) {
            // Close modal
            const modal = bootstrap.Modal.getInstance(document.getElementById('editDeviceModal'));
            if (modal) modal.hide();

            // Clear form
            document.getElementById('editDeviceForm').reset();
            currentEditDeviceId = null;

            // Show success alert
            showAlert(`Device updated successfully`, 'success');

            // Reload page after short delay
            setTimeout(() => {
                window.location.reload();
            }, 1500);
        } else {
            showAlert(result.error || 'Failed to update device', 'danger');
        }
    } catch (error) {
        console.error('Error updating device:', error);
        showAlert('Error updating device. Please try again.', 'danger');
    }
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

// Function untuk confirm delete
async function confirmDelete() {
    console.log('confirmDelete called');
    
    if (!currentDeleteDeviceId) {
        showAlert('No device selected', 'danger');
        return;
    }

    try {
        const response = await fetch(`/api/devices/${currentDeleteDeviceId}`, {
            method: 'DELETE'
        });

        const result = await response.json();

        if (response.ok && result.success) {
            // Close modal
            const modal = bootstrap.Modal.getInstance(document.getElementById('deleteDeviceModal'));
            if (modal) modal.hide();

            currentDeleteDeviceId = null;

            // Show success alert
            showAlert(`Device deleted successfully`, 'success');

            // Reload page after short delay
            setTimeout(() => {
                window.location.reload();
            }, 1500);
        } else {
            showAlert(result.error || 'Failed to delete device', 'danger');
        }
    } catch (error) {
        console.error('Error deleting device:', error);
        showAlert('Error deleting device. Please try again.', 'danger');
    }
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    console.log('devices.js initialized');

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
    document.getElementById('addDeviceModal')?.addEventListener('hidden.bs.modal', function() {
        document.getElementById('addDeviceForm').reset();
    });

    document.getElementById('editDeviceModal')?.addEventListener('hidden.bs.modal', function() {
        document.getElementById('editDeviceForm').reset();
        currentEditDeviceId = null;
    });

    document.getElementById('deleteDeviceModal')?.addEventListener('hidden.bs.modal', function() {
        currentDeleteDeviceId = null;
    });

    console.log('devices.js ready');
});

// Export functions untuk global use
window.saveDevice = saveDevice;
window.updateDevice = updateDevice;
window.openEditModal = openEditModal;
window.openDeleteModal = openDeleteModal;
window.confirmDelete = confirmDelete;

console.log('devices.js functions exported to window');