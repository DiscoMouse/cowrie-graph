// static/js/attacks-by-day.js

const toggleSwitch = document.querySelector('.theme-switch input[type="checkbox"]');
const chartDiv = document.getElementById('chart');
let chart;

function renderChart(theme) {
    if (chart) {
        chart.dispose();
    }
    chart = echarts.init(chartDiv, theme);
    chart.showLoading();

    fetch('/api/v1/attacks-by-day')
        .then(response => response.json())
        .then(data => {
            chart.hideLoading();
            const dates = data.map(item => item.date);
            const successes = data.map(item => item.successes);
            const failures = data.map(item => item.failures);
            const option = { 
                title: { text: 'Attacks By Day', textStyle: { color: theme === 'dark' ? '#fff' : '#333' } }, 
                tooltip: { trigger: 'axis' }, 
                legend: { data: ['Failures', 'Successes'], textStyle: { color: theme === 'dark' ? '#fff' : '#333' } }, 
                xAxis: { type: 'category', boundaryGap: false, data: dates }, 
                yAxis: { type: 'value' }, 
                dataZoom: [{ type: 'slider', start: 70, end: 100 }], 
                series: [ 
                    { name: 'Failures', type: 'line', stack: 'Total', areaStyle: {}, emphasis: { focus: 'series' }, data: failures }, 
                    { name: 'Successes', type: 'line', stack: 'Total', areaStyle: {}, emphasis: { focus: 'series' }, data: successes } 
                ] 
            };
            chart.setOption(option);
        });
}

function setTheme(theme) {
    if (theme === 'dark') {
        document.body.classList.add('dark-mode');
        toggleSwitch.checked = true;
        renderChart('dark');
    } else {
        document.body.classList.remove('dark-mode');
        toggleSwitch.checked = false;
        renderChart(null);
    }
}

const currentTheme = localStorage.getItem('theme') || 'dark';
setTheme(currentTheme);

toggleSwitch.addEventListener('change', function(e) {
    const theme = e.target.checked ? 'dark' : 'light';
    localStorage.setItem('theme', theme);
    setTheme(theme);
});

window.addEventListener('resize', function() {
    if (chart) {
        chart.resize();
    }
});