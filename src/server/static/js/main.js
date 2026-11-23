// WordList main.js
// Theme handling
(function() {
    let savedTheme = 'auto';
    try { savedTheme = localStorage.getItem('theme') || 'auto'; } catch (_) {}
    const media = window.matchMedia('(prefers-color-scheme: light)');
    const apply = () => document.documentElement.setAttribute('data-theme', savedTheme === 'auto' ? (media.matches ? 'light' : 'dark') : savedTheme);
    apply();
    if (media.addEventListener) media.addEventListener('change', apply);
    const toggle = document.getElementById('theme-toggle');
    if (toggle) toggle.addEventListener('click', () => {
        savedTheme = savedTheme === 'auto' ? 'light' : savedTheme === 'light' ? 'dark' : 'auto';
        try { localStorage.setItem('theme', savedTheme); } catch (_) {}
        toggle.textContent = `Theme: ${savedTheme}`;
        apply();
    });
    if (toggle) toggle.textContent = `Theme: ${savedTheme}`;
})();

// Copy to clipboard utility
function copyToClipboard(text) {
    return navigator.clipboard.writeText(text);
}

// Format numbers with commas
function formatNumber(num) {
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',');
}

if ('serviceWorker' in navigator) {
    window.addEventListener('load', () => navigator.serviceWorker.register('/sw.js').catch(() => {}));
}
