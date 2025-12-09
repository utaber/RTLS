// devices.js - Handle CRUD operations with modals

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
    
    // Scroll to top
    window.scrollTo({ top: 0, behavior: 'smooth' });
    
    // Auto dismiss after 5 seconds
    setTimeout(() => {
        const alert = alertContainer.querySelector('.alert');
        if (alert) {
            alert.remove();
        }
    }, 5000);
}

// Save new device
async function saveDevice() {
    const formData = new FormData();
    formData.append('name', document.getElementById('add_name').value);
    formData.append('device_id', document.getElementById('add_device_id').value);
    formData.append('type', document.getElementById('add_type').value);
    formData.append('status', document.getElementById('add_status').value);
    formData.append('location', document.getElementById('add_location').value);
    
    try {
        const response = await fetch('/devices/add', {
            method: 'POST',
            body: formData
        });
        
        if (response.redirected) {
            // Close modal
            const modal = bootstrap.Modal.getInstance(document.getElementById('addDeviceModal'));
            modal.hide();
            
            // Reload page to show new device
            window.location.href = response.url;
        } else {
            showAlert('danger', 'Failed to add device');
        }
    } catch (error) {
        console.error('Error:', error);
        showAlert('danger', 'An error occurred while adding device');
    }
}

// Open edit modal and load device data
async function openEditModal(id) {
    try {
        const response = await fetch(`/api/devices/${id}`);
        const result = await response.json();
        
        if (result.success) {
            const device = result.device;
            
            // Fill form
            document.getElementById('edit_id').value = device.ID;
            document.getElementById('edit_name').value = device.Name;
            document.getElementById('edit_device_id').value = device.DeviceID;
            document.getElementById('edit_type').value = device.Type;
            document.getElementById('edit_status').value = device.Status;
            document.getElementById('edit_location').value = device.Location;
            
            // Show modal
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

// Update device
async function updateDevice() {
    const id = document.getElementById('edit_id').value;
    const formData = new FormData();
    formData.append('name', document.getElementById('edit_name').value);
    formData.append('device_id', document.getElementById('edit_device_id').value);
    formData.append('type', document.getElementById('edit_type').value);
    formData.append('status', document.getElementById('edit_status').value);
    formData.append('location', document.getElementById('edit_location').value);
    
    try {
        const response = await fetch(`/devices/${id}/edit`, {
            method: 'POST',
            body: formData
        });
        
        if (response.redirected) {
            // Close modal
            const modal = bootstrap.Modal.getInstance(document.getElementById('editDeviceModal'));
            modal.hide();
            
            // Reload page
            window.location.href = response.url;
        } else {
            showAlert('danger', 'Failed to update device');
        }
    } catch (error) {
        console.error('Error:', error);
        showAlert('danger', 'An error occurred while updating device');
    }
}

// Open delete modal and load device info
async function openDeleteModal(id) {
    try {
        const response = await fetch(`/api/devices/${id}`);
        const result = await response.json();
        
        if (result.success) {
            const device = result.device;
            
            // Fill info
            document.getElementById('delete_id').value = device.ID;
            document.getElementById('delete_name').textContent = device.Name;
            document.getElementById('delete_device_id').textContent = device.DeviceID;
            
            // Show modal
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

// Confirm delete
async function confirmDelete() {
    const id = document.getElementById('delete_id').value;
    
    try {
        const response = await fetch(`/devices/${id}/delete`, {
            method: 'POST'
        });
        
        if (response.redirected) {
            // Close modal
            const modal = bootstrap.Modal.getInstance(document.getElementById('deleteDeviceModal'));
            modal.hide();
            
            // Reload page
            window.location.href = response.url;
        } else {
            showAlert('danger', 'Failed to delete device');
        }
    } catch (error) {
        console.error('Error:', error);
        showAlert('danger', 'An error occurred while deleting device');
    }
}

// Reset form when modal is hidden
document.getElementById('addDeviceModal').addEventListener('hidden.bs.modal', function () {
    document.getElementById('addDeviceForm').reset();
});

document.getElementById('editDeviceModal').addEventListener('hidden.bs.modal', function () {
    document.getElementById('editDeviceForm').reset();
});