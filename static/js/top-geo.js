// static/js/top-geo.js

const toggleSwitch = document.querySelector('.theme-switch input[type="checkbox"]');
let currentTheme = localStorage.getItem('theme') || 'dark';

const charts = {
    topCountries: echarts.init(document.getElementById('topCountriesChart'), currentTheme),
    topCities: echarts.init(document.getElementById('topCitiesChart'), currentTheme),
    topOrgs: echarts.init(document.getElementById('topOrgsChart'), currentTheme)
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
          hideOverlap: false
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

    renderBarChart(charts.topCountries, '/api/v1/top-countries', data => ({
        labels: data.map(i => i.country_code),
        values: data.map(i => i.count)
    }), 'Top Countries');

    renderBarChart(charts.topCities, '/api/v1/top-cities', data => ({
        labels: data.map(i => i.city),
        values: data.map(i => i.count)
    }), 'Top Cities');

    renderBarChart(charts.topOrgs, '/api/v1/top-orgs', data => ({
        labels: data.map(i => i.organization),
        values: data.map(i => i.count)
    }), 'Top ISPs/Orgs');
}

function setTheme(theme) {
      if (theme === 'dark') { document.body.classList.add('dark-mode'); toggleSwitch.checked = true; } 
      else { document.body.classList.remove('dark-mode'); toggleSwitch.checked = false; }
      
      Object.values(charts).forEach(chart => chart.dispose());
      charts.topCountries = echarts.init(document.getElementById('topCountriesChart'), theme);
      charts.topCities = echarts.init(document.getElementById('topCitiesChart'), theme);
      charts.topOrgs = echarts.init(document.getElementById('topOrgsChart'), theme);
      renderCharts(theme);
}

setTheme(currentTheme);
toggleSwitch.addEventListener('change', e => { currentTheme = e.target.checked ? 'dark' : 'light'; localStorage.setItem('theme', currentTheme); setTheme(currentTheme); });
window.addEventListener('resize', () => Object.values(charts).forEach(chart => chart.resize()));