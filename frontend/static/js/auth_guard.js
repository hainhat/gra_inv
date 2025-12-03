// Ẩn nội dung trang khi đang kiểm tra
document.addEventListener("DOMContentLoaded", () => {
    // Thêm overlay loading
    const loadingOverlay = document.createElement("div");
    loadingOverlay.id = "auth-loading";
    loadingOverlay.innerHTML = `
    <div style="
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      background: rgba(255, 255, 255, 0.95);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 9999;
    ">
      <div style="text-align: center;">
        <div style="
          border: 3px solid #f3f3f3;
          border-top: 3px solid #3498db;
          border-radius: 50%;
          width: 40px;
          height: 40px;
          animation: spin 1s linear infinite;
          margin: 0 auto 16px;
        "></div>
        <p style="color: #666; font-size: 14px;">Đang kiểm tra...</p>
      </div>
    </div>
    <style>
      @keyframes spin {
        0% { transform: rotate(0deg); }
        100% { transform: rotate(360deg); }
      }
    </style>
  `;
    document.body.appendChild(loadingOverlay);

    // Kiểm tra auth
    checkAuthAndRedirect();
});

/**
 * Kiểm tra authentication và redirect nếu cần
 */
async function checkAuthAndRedirect() {
    const token = apiClient.getAccessToken();
    const loadingOverlay = document.getElementById("auth-loading");

    // Nếu không có token -> cho phép truy cập
    if (!token) {
        removeLoadingOverlay();
        return;
    }

    // Nếu có token -> kiểm tra validity
    try {
        const res = await apiClient.get("/me");

        if (res && res.ok) {
            const data = await res.json();

            if (data.success && data.user) {
                // User đã đăng nhập -> redirect về home
                showRedirectNotification(data.user.full_name);

                setTimeout(() => {
                    window.location.href = "/";
                }, 1500);
                return;
            }
        }

        // Token không hợp lệ -> xóa và cho phép truy cập
        apiClient.clearSession();
        removeLoadingOverlay();

    } catch (err) {
        console.error("Error checking auth:", err);
        apiClient.clearSession();
        removeLoadingOverlay();
    }
}

/**
 * Xóa loading overlay
 */
function removeLoadingOverlay() {
    const overlay = document.getElementById("auth-loading");
    if (overlay) {
        overlay.style.opacity = "0";
        overlay.style.transition = "opacity 0.3s ease";
        setTimeout(() => overlay.remove(), 300);
    }
}

/**
 * Hiển thị thông báo redirect
 */
function showRedirectNotification(userName) {
    const overlay = document.getElementById("auth-loading");
    if (overlay) {
        overlay.innerHTML = `
      <div style="
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(255, 255, 255, 0.95);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 9999;
      ">
        <div style="
          background: white;
          padding: 32px;
          border-radius: 12px;
          box-shadow: 0 10px 40px rgba(0,0,0,0.1);
          text-align: center;
          max-width: 400px;
        ">
          <div style="
            width: 64px;
            height: 64px;
            background: #10B981;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            margin: 0 auto 20px;
          ">
            <svg style="width: 32px; height: 32px;" fill="white" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
            </svg>
          </div>
          <h3 style="
            font-size: 20px;
            font-weight: 600;
            color: #111827;
            margin-bottom: 8px;
          ">Xin chào, ${escapeHtml(userName)}!</h3>
          <p style="
            color: #6B7280;
            font-size: 14px;
            margin-bottom: 16px;
          ">Bạn đã đăng nhập rồi. Đang chuyển về trang chủ...</p>
          <div style="
            width: 200px;
            height: 4px;
            background: #E5E7EB;
            border-radius: 2px;
            margin: 0 auto;
            overflow: hidden;
          ">
            <div style="
              width: 100%;
              height: 100%;
              background: #3B82F6;
              animation: progress 1.5s ease;
            "></div>
          </div>
        </div>
      </div>
      <style>
        @keyframes progress {
          from { transform: translateX(-100%); }
          to { transform: translateX(0); }
        }
      </style>
    `;
    }
}

/**
 * Escape HTML để tránh XSS
 */
function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return String(text).replace(/[&<>"']/g, m => map[m]);
}