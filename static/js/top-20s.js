// static/js/top-20s.js

const toggleSwitch = document.querySelector('.theme-switch input[type="checkbox"]');
let currentTheme = localStorage.getItem('theme') || 'dark';

const charts = {
    topPasswords: echarts.init(document.getElementById('topPasswordsChart'), currentTheme),
    topUsernames: echarts.init(document.getElementById('topUsernamesChart'), currentTheme),
    topIPs: echarts.init(document.getElementById('topIPsChart'), currentTheme),
    topClients: echarts.init(document.getElementById('topClientsChart'), currentTheme)
};

function createBarChartOption(title, yAxisData, seriesData, seriesName, theme) {
  const maxValue = Math.max(...seriesData);
  const axisMax = Math.ceil(maxValue * 1.05);
  
  return {
    animationDuration: 2000, // <-- animation speed
    title: { text: title, textStyle: { color: theme === 'dark' ? '#fff' : '#333' } },
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    grid: { left: '25%', containLabel: true },
    xAxis: { type: 'value', max: axisMax },
    yAxis: { 
      type: 'category', 
      data: yAxisData.reverse(),
      axisLabel: {
          interval: 0,
          hideOverlap: false // <-- fix is now included
      }
    },
    series: [{ name: seriesName, type: 'bar', data: seriesData.reverse() }]
  };
}

function renderCharts(theme) {
    const renderBarChart = (chartInstance, endpoint, dataMapper, seriesName) => {
        fetch(endpoint)
            .then(r => r.json())
            .then(data => {
                if (data && data.length > 0) {
                    // Calculate height: 25px per bar + 70px base for title/padding
                    const chartHeight = (data.length * 25) + 70;
                    chartInstance.getDom().style.height = chartHeight + 'px';
                    chartInstance.resize();

                    const { labels, values } = dataMapper(data);
                    chartInstance.setOption(createBarChartOption(seriesName, labels, values, seriesName, theme));
                }
            });
    };

    renderBarChart(charts.topPasswords, '/api/v1/top-passwords', data => ({
        labels: data.map(i => i.password),
        values: data.map(i => i.count)
    }), 'Top Passwords');

    renderBarChart(charts.topUsernames, '/api/v1/top-usernames', data => ({
        labels: data.map(i => i.username),
        values: data.map(i => i.count)
    }), 'Top Usernames');

    renderBarChart(charts.topIPs, '/api/v1/top-ips', data => ({
        labels: data.map(i => i.ip),
        values: data.map(i => i.count)
    }), 'Top Attacker IPs');

    renderBarChart(charts.topClients, '/api/v1/top-clients', data => ({
        labels: data.map(i => i.version),
        values: data.map(i => i.count)
    }), 'Top SSH Clients');
}

function setTheme(theme) {
      if (theme === 'dark') { document.body.classList.add('dark-mode'); toggleSwitch.checked = true; } 
      else { document.body.classList.remove('dark-mode'); toggleSwitch.checked = false; }
      
      Object.values(charts).forEach(chart => chart.dispose());
      charts.topPasswords = echarts.init(document.getElementById('topPasswordsChart'), theme);
      charts.topUsernames = echarts.init(document.getElementById('topUsernamesChart'), theme);
      charts.topIPs = echarts.init(document.getElementById('topIPsChart'), theme);
      charts.topClients = echarts.init(document.getElementById('topClientsChart'), theme);
      renderCharts(theme);
}

setTheme(currentTheme);
toggleSwitch.addEventListener('change', e => { currentTheme = e.target.checked ? 'dark' : 'light'; localStorage.setItem('theme', currentTheme); setTheme(currentTheme); });
window.addEventListener('resize', () => Object.values(charts).forEach(chart => chart.resize()));