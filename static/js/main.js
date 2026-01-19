// Minimal JavaScript - progressive enhancement only
(function() {
    'use strict';

    // Page load complete
    document.documentElement.classList.add('js-loaded');

    // Copy email to clipboard
    const copyableCards = document.querySelectorAll('.contact-card-copyable');

    copyableCards.forEach(function(card) {
        card.addEventListener('click', function() {
            const email = this.dataset.email;
            if (!email) return;

            navigator.clipboard.writeText(email).then(function() {
                showToast('Email copied to clipboard');
            }).catch(function() {
                // Fallback for older browsers
                const textarea = document.createElement('textarea');
                textarea.value = email;
                textarea.style.position = 'fixed';
                textarea.style.opacity = '0';
                document.body.appendChild(textarea);
                textarea.select();
                document.execCommand('copy');
                document.body.removeChild(textarea);
                showToast('Email copied to clipboard');
            });
        });
    });

    function showToast(message) {
        // Remove existing toast if any
        const existing = document.querySelector('.toast');
        if (existing) existing.remove();

        const toast = document.createElement('div');
        toast.className = 'toast';
        toast.textContent = message;
        document.body.appendChild(toast);

        // Trigger animation
        requestAnimationFrame(function() {
            toast.classList.add('toast-visible');
        });

        // Remove after delay
        setTimeout(function() {
            toast.classList.remove('toast-visible');
            setTimeout(function() {
                toast.remove();
            }, 300);
        }, 2500);
    }
})();
