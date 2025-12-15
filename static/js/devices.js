// devices.js - Handle CRUD operations with modals + REST API (Firebase)

// Show alert message
function showAlert(type, message) {
    const alertContainer = document.getElementById('alertContainer');
    const alertHtml = `
        <div class="alert alert-${type} alert-dismissible fade show" role="alert">
            <i class="bi bi-${type === 'success' ? 'check-circle' : 'exclamation-triangle'} me-2"></i>
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        </div>
    `;
    
    alertContainer.innerHTML = alertHtml;
    window.scrollTo({ top: 0, behavior: 'smooth' });
    
    setTimeout(() => {
        const alert = alertContainer.querySelector('.alert');
        if (alert) alert.remove();
    }, 5000);
}

// Save new device (REST API)
async function saveDevice() {
    const formData = new FormData();
    formData.append('device_id', document.getElementById('add_device_id').value);
    formData.append('name', document.getElementById('add_name').value);
    formData.append('status', document.getElementById('add_status').value);
    formData.append('latitude', document.getElementById('add_latitude').value || '1.1045');
    formData.append('longitude', document.getElementById('add_longitude').value || '104.0305');
    
    try {
        const response = await fetch('/api/devices', {
            method: 'POST',
            body: formData
        });
        
        const result = await response.json();
        
        if (result.success) {
            const modal = bootstrap.Modal.getInstance(document.getElementById('addDeviceModal'));
            modal.hide();
            showAlert('success', result.message);
            setTimeout(() => location.reload(), 1000);
        } else {
            showAlert('danger', result.error || 'Failed to add device');
        }
    } catch (error) {
        console.error('Error:', error);
        showAlert('danger', 'An error occurred while adding device');
    }
}

// Open edit modal and load device data (REST API)
async function openEditModal(deviceId) {
    try {
        const response = await fetch(`/api/devices/${deviceId}`);
        const result = await response.json();
        
        if (result.success) {
            const device = result.device;
            
            document.getElementById('edit_device_id').value = device.device_id;
            document.getElementById('edit_name').value = device.name;
            document.getElementById('edit_status').value = device.status;
            document.getElementById('edit_latitude').value = device.latitude;
            document.getElementById('edit_longitude').value = device.longitude;
            
            const modal = new bootstrap.Modal(document.getElementById('editDeviceModal'));
            modal.show();
        } else {
            showAlert('danger', 'Failed to load device data');
        }
    } catch (error) {
        console.error('Error:', error);
        showAlert('danger', 'An error occurred while loading device data');
    }
}

// Update device (REST API)
async function updateDevice() {
    const deviceId = document.getElementById('edit_device_id').value;
    const formData = new FormData();
    formData.append('device_id', deviceId);
    formData.append('name', document.getElementById('edit_name').value);
    formData.append('status', document.getElementById('edit_status').value);
    formData.append('latitude', document.getElementById('edit_latitude').value);
    formData.append('longitude', document.getElementById('edit_longitude').value);
    
    try {
        const response = await fetch(`/api/devices/${deviceId}`, {
            method: 'PUT',
            body: formData
        });
        
        const result = await response.json();
        
        if (result.success) {
            const modal = bootstrap.Modal.getInstance(document.getElementById('editDeviceModal'));
            modal.hide();
            showAlert('success', result.message);
            setTimeout(() => location.reload(), 1000);
        } else {
            showAlert('danger', result.error || 'Failed to update device');
        }
    } catch (error) {
        console.error('Error:', error);
        showAlert('danger', 'An error occurred while updating device');
    }
}

// Open delete modal and load device info (REST API)
async function openDeleteModal(deviceId) {
    try {
        const response = await fetch(`/api/devices/${deviceId}`);
        const result = await response.json();
        
        if (result.success) {
            const device = result.device;
            
            document.getElementById('delete_device_id').value = device.device_id;
            document.getElementById('delete_name').textContent = device.name;
            document.getElementById('delete_id').textContent = device.device_id;
            
            const modal = new bootstrap.Modal(document.getElementById('deleteDeviceModal'));
            modal.show();
        } else {
            showAlert('danger', 'Failed to load device data');
        }
    } catch (error) {
        console.error('Error:', error);
        showAlert('danger', 'An error occurred while loading device data');
    }
}

// Confirm delete (REST API)
async function confirmDelete() {
    const deviceId = document.getElementById('delete_device_id').value;
    
    try {
        const response = await fetch(`/api/devices/${deviceId}`, {
            method: 'DELETE'
        });
        
        const result = await response.json();
        
        if (result.success) {
            const modal = bootstrap.Modal.getInstance(document.getElementById('deleteDeviceModal'));
            modal.hide();
            showAlert('success', result.message);
            setTimeout(() => location.reload(), 1000);
        } else {
            showAlert('danger', result.error || 'Failed to delete device');
        }
    } catch (error) {
        console.error('Error:', error);
        showAlert('danger', 'An error occurred while deleting device');
    }
}

// Reset form when modal is hidden
document.addEventListener('DOMContentLoaded', function() {
    const addModal = document.getElementById('addDeviceModal');
    const editModal = document.getElementById('editDeviceModal');
    
    if (addModal) {
        addModal.addEventListener('hidden.bs.modal', function () {
            document.getElementById('addDeviceForm').reset();
        });
    }
    
    if (editModal) {
        editModal.addEventListener('hidden.bs.modal', function () {
            document.getElementById('editDeviceForm').reset();
        });
    }
});