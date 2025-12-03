// 👁️ Hiển thị / ẩn mật khẩu
function toggleVisibility(input, icon) {
    if (!input || !icon) return;
    icon.addEventListener('click', () => {
        const type = input.getAttribute('type') === 'password' ? 'text' : 'password';
        input.setAttribute('type', type);
        icon.setAttribute('fill', type === 'password' ? '#bbb' : '#1d4ed8');
    });
}

// 🔐 Lưu thông tin user sau khi đăng nhập / đăng ký
function saveUserSession(data) {
    apiClient.setSession(data);
}

// ========================
// 🧾 Đăng nhập
// ========================
function setupLoginPage() {
    const loginForm = document.getElementById('loginForm');
    if (!loginForm) return; // không có form => không phải trang login

    const emailInput = loginForm.querySelector('input[name="email"]');
    const passwordInput = document.getElementById('password');
    const togglePassword = document.getElementById('togglePassword');

    toggleVisibility(passwordInput, togglePassword);

    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const email = emailInput.value.trim();
        const password = passwordInput.value.trim();

        if (!email || !password) {
            alert('Vui lòng nhập đầy đủ email và mật khẩu!');
            return;
        }

        try {
            const res = await apiClient.post('/login', { email, password });
            if (!res) return; // Lỗi mạng hoặc gì đó

            const result = await res.json();

            if (result.success) {
                saveUserSession(result);
                window.location.href = '/';
            } else {
                alert(result.message || 'Sai thông tin đăng nhập!');
            }
        } catch (err) {
            console.error('Login error:', err);
            alert('Không thể kết nối tới server.');
        }
    });
}

// ========================
// 🧾 Đăng ký
// ========================
function setupRegisterPage() {
    const form = document.getElementById('registerForm');
    if (!form) return; // không có form => không phải trang register

    const passwordInput = document.getElementById('password');
    const confirmInput = document.getElementById('confirmPassword');
    const passwordError = document.getElementById('passwordError');
    const togglePassword = document.getElementById('togglePassword');
    const toggleConfirm = document.getElementById('toggleConfirm');
    const emailInput = form.querySelector('input[name="email"]');

    toggleVisibility(passwordInput, togglePassword);
    toggleVisibility(confirmInput, toggleConfirm);

    // ✅ Kiểm tra email trùng
    async function checkEmailExists(email) {
        try {
            const res = await apiClient.get(`/check-email?email=${encodeURIComponent(email)}`);
            if (!res) return false;
            const data = await res.json();
            return data.exists === true;
        } catch (err) {
            console.error('Email check error:', err);
            return false;
        }
    }

    emailInput.addEventListener('blur', async () => {
        const email = emailInput.value.trim();
        if (!email) return;
        const exists = await checkEmailExists(email);
        const errorEl = document.getElementById('emailError');
        if (exists) {
            emailInput.classList.add('border-red-500', 'focus:border-red-500');
            emailInput.classList.remove('border-slate-300');
            if (errorEl) errorEl.classList.remove('hidden');
        } else {
            emailInput.classList.remove('border-red-500');
            emailInput.classList.add('border-slate-300');
            if (errorEl) errorEl.classList.add('hidden');
        }
    });

    // 🧾 Gửi form
    form.addEventListener('submit', async (e) => {
        e.preventDefault();

        const full_name = form.querySelector('input[name="full_name"]').value.trim();
        const email = emailInput.value.trim();
        const phone = form.querySelector('input[name="phone"]').value.trim();
        const password = passwordInput.value.trim();
        const confirm_password = confirmInput.value.trim();

        if (password !== confirm_password) {
            passwordError.classList.remove('hidden');
            return;
        } else {
            passwordError.classList.add('hidden');
        }

        const exists = await checkEmailExists(email);
        if (exists) {
            alert('Email này đã tồn tại, vui lòng dùng email khác.');
            return;
        }

        try {
            const res = await apiClient.post('/register', { full_name, email, phone, password });
            if (!res) return;

            const data = await res.json();

            if (data.success) {
                saveUserSession(data);
                alert('🎉 Đăng ký thành công!');
                window.location.href = '/';
            } else {
                alert(data.message || 'Đăng ký thất bại!');
            }
        } catch (error) {
            console.error('Error:', error);
            alert('Không thể kết nối tới server.');
        }
    });
}

// ========================
// 🚀 Khởi chạy tương ứng trang
// ========================
document.addEventListener('DOMContentLoaded', () => {
    setupLoginPage();
    setupRegisterPage();
});