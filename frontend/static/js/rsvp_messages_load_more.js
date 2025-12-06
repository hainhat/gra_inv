class RSVPMessagesManager {
    constructor() {
        this.container = document.getElementById('messagesContainer');
        this.loadingEl = document.getElementById('loadingMessages');
        this.emptyStateEl = document.getElementById('emptyState');
        this.messageCountEl = document.getElementById('messageCount');
        this.loadMoreBtnContainer = document.getElementById('loadMoreContainer');

        // Config
        this.API_BASE_URL = '/api';
        this.PAGE_SIZE = 3;

        // State
        this.currentPage = 0;
        this.isLoading = false;
        this.hasMore = true;

        if (this.container) {
            this.init();
        }
    }

    init() {
        this.createLoadMoreButton();
        this.loadMessages();
    }

    createLoadMoreButton() {
        if (!this.loadMoreBtnContainer) return;

        const button = document.createElement('button');
        button.id = 'loadMoreBtn';
        button.className = 'w-full mt-4 bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-3 px-6 rounded-lg transition-colors';
        button.innerHTML = 'Xem thêm';
        button.style.display = 'none';

        button.addEventListener('click', () => this.loadMessages());

        this.loadMoreBtnContainer.appendChild(button);
    }

    async loadMessages() {
        if (this.isLoading || !this.hasMore) return;

        try {
            this.isLoading = true;
            const nextPage = this.currentPage + 1;

            if (nextPage === 1) {
                this.showLoading(true);
            } else {
                this.setButtonLoading(true);
            }

            const response = await fetch(
                `${this.API_BASE_URL}/rsvp/messages?page=${nextPage}&limit=${this.PAGE_SIZE}`
            );
            const data = await response.json();

            if (data.success) {
                const messages = data.data || [];
                const pagination = data.pagination || {};

                this.currentPage = nextPage;
                this.hasMore = this.currentPage < pagination.totalPages;

                // Update message count
                if (this.messageCountEl) {
                    this.messageCountEl.textContent = `${pagination.total || 0} tin nhắn`;
                }

                if (nextPage === 1 && messages.length === 0) {
                    this.showEmptyState(true);
                    this.hideLoadMoreButton();
                } else {
                    this.showEmptyState(false);
                    this.renderMessages(messages);

                    if (this.hasMore) {
                        this.showLoadMoreButton();
                    } else {
                        this.hideLoadMoreButton();
                    }
                }
            }
        } catch (error) {
            console.error('Error loading messages:', error);
            alert('Không thể tải tin nhắn. Vui lòng thử lại.');
        } finally {
            this.isLoading = false;
            this.showLoading(false);
            this.setButtonLoading(false);
        }
    }

    renderMessages(messages) {
        messages.forEach(msg => {
            const element = this.createMessageElement(msg);
            this.container.appendChild(element);
        });
    }

    createMessageElement(msg) {
        const div = document.createElement('div');
        div.className = 'flex gap-2.5 py-3 border-b border-gray-100 dark:border-gray-700 last:border-0';

        const timeAgo = this.formatTimeAgo(msg.created_at);
        const avatarUrl = msg.avatar || 'https://res.cloudinary.com/dcncfkvwv/image/upload/v1733476463/sum8iqnxhdgdyj6zcc2l.jpg';

        div.innerHTML = `
            <div class="flex-shrink-0">
                <img src="${avatarUrl}"
                     alt="${this.escapeHtml(msg.name)}"
                     class="w-10 h-10 rounded-full object-cover"
                     onerror="this.src='https://res.cloudinary.com/dcncfkvwv/image/upload/v1733476463/sum8iqnxhdgdyj6zcc2l.jpg'">
            </div>
            <div class="flex-1 min-w-0">
                <div class="flex items-baseline justify-between mb-1 gap-2">
                    <h3 class="font-semibold text-gray-900 dark:text-white text-sm truncate">
                        ${this.escapeHtml(msg.name)}
                    </h3>
                    <span class="text-xs text-gray-400 dark:text-gray-500 flex-shrink-0">
                        ${timeAgo}
                    </span>
                </div>
                <p class="text-gray-700 dark:text-gray-300 text-sm break-words">
                    ${this.escapeHtml(msg.message)}
                </p>
            </div>
        `;

        return div;
    }

    showLoading(show) {
        if (this.loadingEl) {
            this.loadingEl.classList.toggle('hidden', !show);
        }
    }

    showLoadMoreButton() {
        const btn = document.getElementById('loadMoreBtn');
        if (btn) btn.style.display = 'block';
    }

    hideLoadMoreButton() {
        const btn = document.getElementById('loadMoreBtn');
        if (btn) btn.style.display = 'none';
    }

    setButtonLoading(loading) {
        const btn = document.getElementById('loadMoreBtn');
        if (btn) {
            btn.disabled = loading;
            btn.textContent = loading ? 'Đang tải...' : 'Xem thêm';
            if (loading) {
                btn.classList.add('opacity-50', 'cursor-not-allowed');
            } else {
                btn.classList.remove('opacity-50', 'cursor-not-allowed');
            }
        }
    }

    showEmptyState(show) {
        if (this.emptyStateEl) {
            this.emptyStateEl.classList.toggle('hidden', !show);
        }
    }

    formatTimeAgo(dateString) {
        const date = new Date(dateString);
        const now = new Date();
        const seconds = Math.floor((now - date) / 1000);

        if (seconds < 60) return 'Vừa xong';
        if (seconds < 3600) return `${Math.floor(seconds / 60)} phút trước`;
        if (seconds < 86400) return `${Math.floor(seconds / 3600)} giờ trước`;
        if (seconds < 604800) return `${Math.floor(seconds / 86400)} ngày trước`;
        if (seconds < 2592000) return `${Math.floor(seconds / 604800)} tuần trước`;
        if (seconds < 31536000) return `${Math.floor(seconds / 2592000)} tháng trước`;
        return `${Math.floor(seconds / 31536000)} năm trước`;
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    reload() {
        this.container.innerHTML = '';
        this.currentPage = 0;
        this.hasMore = true;
        this.hideLoadMoreButton();
        this.loadMessages();
    }
}

// Initialize
let rsvpMessagesManager;

document.addEventListener('DOMContentLoaded', () => {
    rsvpMessagesManager = new RSVPMessagesManager();

    // Reload sau khi submit form
    const rsvpForm = document.getElementById('rsvpForm');
    if (rsvpForm) {
        rsvpForm.addEventListener('submit', () => {
            setTimeout(() => {
                if (rsvpMessagesManager) {
                    rsvpMessagesManager.reload();
                }
            }, 2000);
        });
    }
});
