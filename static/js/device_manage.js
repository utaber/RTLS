// device_manage.js - Handle CRUD operations dengan AJAX

// Form elements
const form = document.getElementById('deviceForm');
const formTitle = document.getElementById('formTitle');
const submitBtn = document.getElementById('submitBtn');
const cancelBtn = document.getElementById('cancelBtn');
const deviceIdInput = document.getElementById('deviceId');

// Form submit handler
form.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const deviceId = deviceIdInput.value;
    const formData = new FormData();
    
    formData.append('name', document.getElementById('name').value);
    formData.append('device_id', document.getElementById('device_id').value);
    formData.append('type', document.getElementById('type').value);
    formData.append('status', document.getElementById('status').value);
    formData.append('location', document.getElementById('location').value);
    
    try {
        let url, method;
        
        if (deviceId) {
            // Update existing device
            url = `/api/devices/${deviceId}`;
            method = 'POST';
        } else {
            // Create new device
            url = '/api/devices';
            method = 'POST';
        }
        
        const response = await fetch(url, {
            method: method,
            body: formData
        });
        
        const result = await response.json();
        
        if (result.success) {
            showAlert('success', result.message);
            form.reset();
            resetForm();
            
            // Reload page to update table
            setTimeout(() => {
                location.reload();
            }, 1000);
        } else {
            showAlert('danger', result.error || 'Failed to save device');
        }
    } catch (error) {
        console.error('Error:', error);
        showAlert('danger', 'An error occurred while saving the device');
    }
});

// Cancel button handler
cancelBtn.addEventListener('click', () => {
    resetForm();
});

// Edit device function
async function editDevice(id) {
    try {
        const response = await fetch(`/api/devices/${id}`);
        const result = await response.json();
        
        if (result.success) {
            const device = result.device;
            
            // Populate form
            deviceIdInput.value = device.ID;
            document.getElementById('name').value = device.Name;
            document.getElementById('device_id').value = device.DeviceID;
            document.getElementById('type').value = device.Type;
            document.getElementById('status').value = device.Status;
            document.getElementById('location').value = device.Location;
            
            // Update form UI
            formTitle.innerHTML = '<i class="bi bi-pencil-square me-2"></i>Edit Device';
            submitBtn.innerHTML = '<i class="bi bi-save me-2"></i>Update Device';
            cancelBtn.style.display = 'block';
            
            // Scroll to form
            window.scrollTo({ top: 0, behavior: 'smooth' });
        } else {
            showAlert('danger', 'Failed to load device data');
        }
    } catch (error) {
        console.error('Error:', error);
        showAlert('danger', 'An error occurred while loading device data');
    }
}

// Delete device function
async function deleteDevice(id) {
    if (!confirm('Are you sure you want to delete this device?')) {
        return;
    }
    
    try {
        const response = await fetch(`/api/devices/${id}`, {
            method: 'DELETE'
        });
        
        const result = await response.json();
        
        if (result.success) {
            showAlert('success', result.message);
            
            // Remove row from table
            const row = document.querySelector(`tr[data-id="${id}"]`);
            if (row) {
                row.remove();
            }
            
            // Reload if no devices left
            const tbody = document.getElementById('deviceTableBody');
            if (tbody.children.length === 0) {
                location.reload();
            }
        } else {
            showAlert('danger', result.error || 'Failed to delete device');
        }
    } catch (error) {
        console.error('Error:', error);
        showAlert('danger', 'An error occurred while deleting the device');
    }
}

// Reset form to default state
function resetForm() {
    form.reset();
    deviceIdInput.value = '';
    formTitle.innerHTML = '<i class="bi bi-plus-circle me-2"></i>Add New Device';
    submitBtn.innerHTML = '<i class="bi bi-save me-2"></i>Save Device';
    cancelBtn.style.display = 'none';
}

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
    
    // Auto dismiss after 5 seconds
    setTimeout(() => {
        const alert = alertContainer.querySelector('.alert');
        if (alert) {
            alert.remove();
        }
    }, 5000);
}

// Make functions globally available
window.editDevice = editDevice;
window.deleteDevice = deleteDevice;