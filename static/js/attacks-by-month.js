// static/js/attacks-by-month.js

const toggleSwitch = document.querySelector('.theme-switch input[type="checkbox"]');
const chartDiv = document.getElementById('chart');
let chart;

function renderChart(theme) {
      if (chart) {
          chart.dispose();
      }
      chart = echarts.init(chartDiv, theme);
      chart.showLoading();

      fetch('/api/v1/attacks-by-month')
          .then(response => response.json())
          .then(data => {
              chart.hideLoading();
              const option = {
                  title: { text: 'Attacks By Month', textStyle: { color: theme === 'dark' ? '#fff' : '#333' } },
                  tooltip: { trigger: 'axis' },
                  legend: { data: ['Failures', 'Successes'], textStyle: { color: theme === 'dark' ? '#fff' : '#333' } }, 
                  xAxis: { type: 'category', data: data.map(item => item.date) },
                  yAxis: { type: 'value' },
                  series: [
                      { name: 'Failures', type: 'bar', stack: 'Total', data: data.map(item => item.failures) },
                      { name: 'Successes', type: 'bar', stack: 'Total', data: data.map(item => item.successes) }
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