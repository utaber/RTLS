// main.js - FIXED VERSION (HttpOnly Support)
console.log('main.js loaded successfully');

// ✅ Fungsi getCookie (Tetap disimpan untuk kegunaan lain, misal ambil email)
function getCookie(name) {
    const nameEQ = name + "=";
    const cookies = document.cookie.split(';');
    for(let i = 0; i < cookies.length; i++) {
        const cookie = cookies[i].trim();
        if (cookie.indexOf(nameEQ) === 0) {
            return cookie.substring(nameEQ.length);
        }
    }
    return null;
}

// ✅ Fungsi showAlert untuk menampilkan pesan
function showAlert(message, type = 'success') {
    let alertContainer = document.getElementById('alertContainer');
    if (!alertContainer) {
        alertContainer = document.createElement('div');
        alertContainer.id = 'alertContainer';
        alertContainer.style.position = 'fixed';
        alertContainer.style.top = '20px';
        alertContainer.style.right = '20px';
        alertContainer.style.zIndex = '9999';
        alertContainer.style.maxWidth = '400px';
        document.body.appendChild(alertContainer);
    }

    const alertHTML = `
        <div class="alert alert-${type} alert-dismissible fade show" role="alert" style="min-width: 300px;">
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

// ✅ Fungsi utama fetchAPI (DIPERBAIKI)
async function fetchAPI(endpoint, method = 'GET', body = null) {
    console.log(`fetchAPI: ${method} ${endpoint}`);
    
    // PERBAIKAN: Kita TIDAK mengecek token di sini karena cookie HttpOnly tidak terlihat oleh JS.
    // Browser akan otomatis mengirim cookie 'auth_token' karena opsi credentials: 'include'.

    // Setup options
    const options = {
        method: method,
        headers: {
            'Content-Type': 'application/json'
            // PERBAIKAN: Tidak perlu header Authorization manual
        },
        credentials: 'include' // PENTING: Ini yang mengirim cookie HttpOnly ke server
    };

    // Tambahkan body jika ada
    if (body) {
        options.body = JSON.stringify(body);
        console.log('Request body:', body);
    }

    try {
        // Eksekusi request
        const response = await fetch(endpoint, options);
        console.log(`Response status: ${response.status}`);
        
        // Handle unauthorized (Jika Server Go merespon 401 karena token invalid/expired)
        if (response.status === 401) {
            console.warn('Unauthorized response from server');
            showAlert('Session expired. Please login again.', 'danger');
            setTimeout(() => {
                window.location.href = '/login';
            }, 2000);
            return response;
        }
        
        return response;
    } catch (error) {
        console.error(`API Error [${method} ${endpoint}]:`, error);
        showAlert(`Network error: ${error.message}`, 'danger');
        throw error;
    }
}

// ✅ Export fungsi ke window object untuk akses global
window.getCookie = getCookie;
window.showAlert = showAlert;
window.fetchAPI = fetchAPI;

console.log('main.js functions exported to window');