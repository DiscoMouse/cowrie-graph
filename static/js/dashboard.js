// static/js/dashboard.js

const toggleSwitch = document.querySelector('.theme-switch input[type="checkbox"]');
const currentTheme = localStorage.getItem('theme') || 'dark';

function setTheme(theme) {
    if (theme === 'dark') {
        document.body.classList.add('dark-mode');
        toggleSwitch.checked = true;
    } else {
        document.body.classList.remove('dark-mode');
        toggleSwitch.checked = false;
    }
}

setTheme(currentTheme);

toggleSwitch.addEventListener('change', function(e) {
    let theme = 'light';
    if (e.target.checked) {
        theme = 'dark';
    }
    document.body.classList.toggle('dark-mode');
    localStorage.setItem('theme', theme);
});