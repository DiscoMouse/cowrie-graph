// static/js/world-map.js

const toggleSwitch = document.querySelector('.theme-switch input[type="checkbox"]');
const chartDiv = document.getElementById('chart');
let chart;

function renderChart(theme) {
      if (chart) {
          chart.dispose();
      }
      chart = echarts.init(chartDiv, theme);
      chart.showLoading();

      fetch('/api/v1/attacks-by-location')
          .then(response => {
              if (!response.ok) {
                  throw new Error(`HTTP error! status: ${response.status}`);
              }
              return response.json();
          })
          .then(data => {
              chart.hideLoading();

              if (!data || data.length === 0) {
                  chart.setOption({ title: { text: 'No location data available', left: 'center', top: 'center', textStyle: { color: theme === 'dark' ? '#fff' : '#333' } } });
                  return;
              }

              const mapData = data.map(item => ({
                  name: item.city || item.country_code,
                  value: [item.longitude, item.latitude, item.count]
              }));

              const option = {
                  title: { text: 'Attacker Origins', textStyle: { color: theme === 'dark' ? '#fff' : '#333' } },
                  tooltip: {
                      trigger: 'item',
                      formatter: function (params) {
                          if (!params.value) return params.name;
                          return `${params.name}<br/>Sessions: ${params.value[2]}`;
                      }
                  },
                  visualMap: {
                      min: 1,
                      max: Math.max(1, ...mapData.map(item => item.value[2])),
                      left: 'left',
                      top: 'bottom',
                      text: ['High', 'Low'],
                      calculable: true,
                      inRange: { color: ['#5291FF', '#E06343', '#B80909'] },
                      textStyle: { color: theme === 'dark' ? '#fff' : '#333' }
                  },
                  geo: {
                      map: 'world',
                      roam: true,
                      emphasis: { label: { show: false }, itemStyle: { areaColor: '#a9b3c4' } },
                      itemStyle: { areaColor: theme === 'dark' ? '#323c48' : '#eee', borderColor: '#111' },
                  },
                  series: [
                      {
                          name: 'Attacks',
                          type: 'scatter',
                          coordinateSystem: 'geo',
                          data: mapData,
                          symbolSize: function (val) {
                              return val ? Math.max(5, Math.log2(val[2]) * 3) : 0;
                          },
                          encode: { value: 2 },
                      }
                  ]
              };
              chart.setOption(option);
          })
          .catch(error => {
              chart.hideLoading();
              console.error("Error fetching or processing data:", error);
              chart.setOption({ title: { text: 'Failed to load map data', subtext: 'Check console for details.', left: 'center', top: 'center', textStyle: { color: theme === 'dark' ? '#fff' : '#333' } } });
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