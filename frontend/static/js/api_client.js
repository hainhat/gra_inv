window.API_URL = 'http://localhost:8080/api';
const API_URL = window.API_URL;

const apiClient = {
    // Helper để lấy token
    getAccessToken: () => localStorage.getItem('access_token'),
    getRefreshToken: () => localStorage.getItem('refresh_token'),
    getUserName: () => localStorage.getItem('user_name'),
    getUserRole: () => localStorage.getItem('user_role'),

    // Helper để lưu session
    setSession: (data) => {
        if (data.access_token) localStorage.setItem('access_token', data.access_token);
        if (data.refresh_token) localStorage.setItem('refresh_token', data.refresh_token);
        if (data.user) {
            localStorage.setItem('user_name', data.user.full_name);
            localStorage.setItem('user_role', data.user.role);
        }
    },

    // Helper để xóa session
    clearSession: () => {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user_name');
        localStorage.removeItem('user_role');
        // Xóa token cũ nếu còn sót lại
        localStorage.removeItem('token');
    },

    // Hàm fetch chính
    request: async (endpoint, options = {}) => {
        let url = endpoint.startsWith('http') ? endpoint : `${API_URL}${endpoint}`;
        
        // Mặc định headers
        const headers = {
            'Content-Type': 'application/json',
            ...options.headers
        };

        // Tự động gắn Access Token nếu có
        const token = apiClient.getAccessToken();
        if (token) {
            headers['Authorization'] = `Bearer ${token}`;
        }

        const config = {
            ...options,
            headers
        };

        try {
            let response = await fetch(url, config);

            // Nếu gặp lỗi 401 (Unauthorized) -> Thử refresh token
            if (response.status === 401) {
                console.warn('⚠️ Access Token hết hạn, đang thử refresh...');
                const refreshSuccess = await apiClient.refreshToken();

                if (refreshSuccess) {
                    // Refresh thành công -> Gọi lại request ban đầu với token mới
                    const newToken = apiClient.getAccessToken();
                    config.headers['Authorization'] = `Bearer ${newToken}`;
                    response = await fetch(url, config);
                } else {
                    // Refresh thất bại -> Đăng xuất
                    console.error('❌ Refresh Token cũng hết hạn hoặc không hợp lệ.');
                    apiClient.clearSession();
                    window.location.href = '/login';
                    return null;
                }
            }

            return response;
        } catch (error) {
            console.error('API Request Error:', error);
            throw error;
        }
    },

    // Hàm refresh token
    refreshToken: async () => {
        const refreshToken = apiClient.getRefreshToken();
        if (!refreshToken) return false;

        try {
            const res = await fetch(`${API_URL}/refresh`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ refresh_token: refreshToken })
            });

            const data = await res.json();

            if (data.success && data.access_token) {
                console.log('✅ Refresh Token thành công!');
                localStorage.setItem('access_token', data.access_token);
                return true;
            } else {
                return false;
            }
        } catch (error) {
            console.error('Refresh Token Error:', error);
            return false;
        }
    },

    // Các method tiện ích
    get: (endpoint) => apiClient.request(endpoint, { method: 'GET' }),
    post: (endpoint, body) => apiClient.request(endpoint, { method: 'POST', body: JSON.stringify(body) }),
    put: (endpoint, body) => apiClient.request(endpoint, { method: 'PUT', body: JSON.stringify(body) }),
    delete: (endpoint) => apiClient.request(endpoint, { method: 'DELETE' })
};

// Expose ra global để dùng ở các file khác
window.apiClient = apiClient;
