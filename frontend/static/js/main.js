document.addEventListener('DOMContentLoaded', () => {
    const API_URL = window.API_URL;
    console.log('Stats loader initialized, API_URL:', API_URL);

    // Load Stats
    async function loadStats() {
        try {
            console.log('Fetching stats from:', `${API_URL}/rsvp/stats`);
            const response = await fetch(`${API_URL}/rsvp/stats`);
            const result = await response.json();
            console.log('Stats response:', result);

            if (result.success) {
                const stats = result.data;
                console.log('Stats data:', stats);
                
                // Update individual stat elements
                const totalEl = document.getElementById('totalCount');
                const yesEl = document.getElementById('yesCount');
                const noEl = document.getElementById('noCount');
                const maybeEl = document.getElementById('maybeCount');
                
                console.log('Elements found:', {
                    total: !!totalEl,
                    yes: !!yesEl,
                    no: !!noEl,
                    maybe: !!maybeEl
                });
                
                if (totalEl) {
                    totalEl.textContent = stats.total || 0;
                    console.log('Set total to:', stats.total);
                }
                if (yesEl) {
                    yesEl.textContent = stats.yes || 0;
                    console.log('Set yes to:', stats.yes);
                }
                if (noEl) {
                    noEl.textContent = stats.no || 0;
                    console.log('Set no to:', stats.no);
                }
                if (maybeEl) {
                    maybeEl.textContent = stats.maybe || 0;
                    console.log('Set maybe to:', stats.maybe);
                }
            }
        } catch (error) {
            console.error('Error loading stats:', error);
        }
    }

    // Expose loadStats globally
    window.loadRSVPStats = loadStats;

    // Initialize
    loadStats();
});