// static/js/country-race.js

const chartDiv = document.getElementById('chart');
const toggleSwitch = document.querySelector('.theme-switch input[type="checkbox"]');
let chart;
let currentTheme = localStorage.getItem('theme') || 'dark';
let animationTimeout;

// Helper to turn a country code like "US" into a flag emoji 🇺🇸
function getFlagEmoji(countryCode) {
    if (!countryCode || countryCode.length !== 2) return countryCode;
    const codePoints = countryCode
      .toUpperCase()
      .split('')
      .map(char =>  127397 + char.charCodeAt(0));
    return String.fromCodePoint(...codePoints);
}

function run(rawData, theme) {
    if (chart) chart.dispose();
    if (animationTimeout) clearTimeout(animationTimeout);

    chart = echarts.init(chartDiv, theme);

    if (!rawData || rawData.length === 0) {
        chart.setOption({ title: { text: 'No country data available yet', left: 'center', top: 'center', textStyle: { color: theme === 'dark' ? '#fff' : '#333' } } });
        return;
    }

    const hours = [...new Set(rawData.map(item => item.hour))].sort();
    const topN = 10;

    const option = { /* ... ECharts options ... */ };
    chart.setOption({
        grid: { top: 10, bottom: 30, left: 150, right: 80 },
        xAxis: { max: 'dataMax', axisLabel: { show: false } },
        yAxis: { type: 'category', inverse: true, max: topN - 1,
            axisLabel: { show: true, textStyle: { fontSize: 16 } }, // Bigger font for emojis
        },
        series: [{ type: 'bar', realtimeSort: true,
            label: { show: true, position: 'right', valueAnimation: true }
        }],
        animationDuration: 0,
        animationDurationUpdate: 2000,
        animationEasing: 'linear',
        animationEasingUpdate: 'linear',
    });
    
    let cumulativeData = {};
    let hourIndex = 0;

    function updateHour() {
        if (hourIndex >= hours.length) return;

        const hour = hours[hourIndex];
        const currentHourData = rawData.filter(item => item.hour === hour);
        
        currentHourData.forEach(item => {
            cumulativeData[item.country_code] = (cumulativeData[item.country_code] || 0) + item.count;
        });

        const sortedData = Object.entries(cumulativeData)
            .sort((a, b) => b[1] - a[1])
            .slice(0, topN);

        chart.setOption({
            graphic: {
                elements: [{
                    type: 'text', right: 160, bottom: 60,
                    style: { text: hour, font: 'bolder 50px monospace', fill: theme === 'dark' ? 'rgba(255, 255, 255, 0.25)' : 'rgba(0, 0, 0, 0.25)' },
                    z: 100
                }]
            },
            yAxis: { data: sortedData.map(item => `${getFlagEmoji(item[0])} ${item[0]}`) },
            series: [{ data: sortedData.map(item => item[1]) }]
        });

        hourIndex++;
        animationTimeout = setTimeout(updateHour, 2000);
    }

    updateHour();
}

fetch('/api/v1/country-race-data')
    .then(response => response.json())
    .then(data => {
        function setTheme(theme) {
            if (theme === 'dark') { document.body.classList.add('dark-mode'); toggleSwitch.checked = true; } 
            else { document.body.classList.remove('dark-mode'); toggleSwitch.checked = false; }
            run(data, theme);
        }
        setTheme(currentTheme);
        toggleSwitch.addEventListener('change', e => {
            currentTheme = e.target.checked ? 'dark' : 'light';
            localStorage.setItem('theme', currentTheme);
            setTheme(currentTheme);
        });
    })
    .catch(error => {
        console.error("Error fetching country bar race data:", error);
        if(chart) chart.dispose();
        chart = echarts.init(chartDiv, currentTheme);
        chart.setOption({ title: { text: 'Failed to load data', subtext: 'Check console for details.', left: 'center', top: 'center', textStyle: { color: currentTheme === 'dark' ? '#fff' : '#333' } } });
    });

window.addEventListener('resize', () => { if (chart) chart.resize(); });