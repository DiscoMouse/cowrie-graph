// static/js/bar-race.js

const chartDiv = document.getElementById('chart');
const toggleSwitch = document.querySelector('.theme-switch input[type="checkbox"]');
let chart;
let currentTheme = localStorage.getItem('theme') || 'dark';
let animationTimeout;

function run(rawData, theme) {
    if (chart) chart.dispose();
    if (animationTimeout) clearTimeout(animationTimeout);

    chart = echarts.init(chartDiv, theme);

    if (!rawData || rawData.length === 0) {
        chart.setOption({ title: { text: 'No session data available yet', left: 'center', top: 'center', textStyle: { color: theme === 'dark' ? '#fff' : '#333' } } });
        return;
    }

    // --- NEW: More robust timeline generation ---
    const dataByHour = new Map();
    rawData.forEach(item => {
        if (!dataByHour.has(item.hour)) {
            dataByHour.set(item.hour, []);
        }
        dataByHour.get(item.hour).push({ ip: item.ip, count: item.count });
    });

    const sortedHours = [...dataByHour.keys()].sort();
    const startTime = new Date(sortedHours[0]);
    const endTime = new Date(sortedHours[sortedHours.length - 1]);
    
    const timeline = [];
    let currentTime = new Date(startTime);
    while (currentTime <= endTime) {
        timeline.push(currentTime.toISOString().slice(0, 16).replace('T', ' '));
        currentTime.setHours(currentTime.getHours() + 1);
    }
    console.log("Generated timeline with " + timeline.length + " hours.");
    // --- END OF NEW LOGIC ---

    const topN = 10;
    const option = { /* ... ECharts options ... */ }; // Placeholder
    chart.setOption({
        grid: { top: 10, bottom: 30, left: 150, right: 80 },
        xAxis: { max: 'dataMax', axisLabel: { show: false } },
        yAxis: { type: 'category', inverse: true, max: topN - 1,
            axisLabel: { show: true, textStyle: { fontSize: 12 } },
        },
        series: [{ type: 'bar', realtimeSort: true,
            label: { show: true, position: 'right', valueAnimation: true }
        }],
        animationDuration: 0,
        animationDurationUpdate: 500,
        animationEasing: 'linear',
        animationEasingUpdate: 'linear',
    });
    
    let cumulativeData = {};
    let hourIndex = 0;

    function updateHour() {
        if (hourIndex >= timeline.length) {
            return; // Animation finished
        }
        const hour = timeline[hourIndex];
        
        if (dataByHour.has(hour)) {
            const currentHourData = dataByHour.get(hour);
            currentHourData.forEach(item => {
                cumulativeData[item.ip] = (cumulativeData[item.ip] || 0) + item.count;
            });
        }
        
        const sortedData = Object.entries(cumulativeData)
            .sort((a, b) => b[1] - a[1])
            .slice(0, topN);
        
        console.log(`Data for hour ${hour}:`, sortedData); // Detailed logging for each frame

        chart.setOption({
            graphic: {
                elements: [{
                    type: 'text', right: 160, bottom: 60,
                    style: {
                        text: hour,
                        font: 'bolder 50px monospace',
                        fill: theme === 'dark' ? 'rgba(255, 255, 255, 0.25)' : 'rgba(0, 0, 0, 0.25)'
                    },
                    z: 100
                }]
            },
            yAxis: { data: sortedData.map(item => item[0]) },
            series: [{ data: sortedData.map(item => item[1]) }]
        });

        hourIndex++;
        animationTimeout = setTimeout(updateHour, 500); // Faster animation
    }

    updateHour();
}

fetch('/api/v1/bar-race-data')
    .then(response => response.json())
    .then(data => {
        console.log("Raw data from API:", data);
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
        console.error("Error fetching bar race data:", error);
        if(chart) chart.dispose();
        chart = echarts.init(chartDiv, currentTheme);
        chart.setOption({ title: { text: 'Failed to load data', subtext: 'Check console for details.', left: 'center', top: 'center', textStyle: { color: currentTheme === 'dark' ? '#fff' : '#333' } } });
    });

window.addEventListener('resize', () => { if (chart) chart.resize(); });