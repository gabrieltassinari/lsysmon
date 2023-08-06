Chart.defaults.color = '#999999'

const memoryCtx = document.getElementById('memory');
const memoryPlot = createChart(memoryCtx, [], 'Free memory (GB/s)');

const swapCtx = document.getElementById('swap');
const swapPlot = createChart(swapCtx, [], 'Swap in use (GB/s)');

const ptimeCtx = document.getElementById('ptime');
const ptimePlot = createChart(ptimeCtx, [], 'Process time');

const cpuCtx = document.getElementById('cpu');
const cpuPlot = createChart(cpuCtx, [], 'Cpu percentage');

cpuPlot.options.maintainAspectRatio = false

ptimePlot.options.maintainAspectRatio = false
ptimePlot.data.datasets = [{
				label: "Utime",
				data: [],
				borderWidth: 1,
				pointRadius: 0,
			},
			{
				label: "Stime",
				data: [],
				borderWidth: 1,
				pointRadius: 0,
			},
			{
				label: "Cutime",
				data: [],
				borderWidth: 1,
				pointRadius: 0,
			},
			{
				label: "Cstime",
				data: [],
				borderWidth: 1,
				pointRadius: 0,
			}]

function createChart(ctx, data, label) {
	return new Chart(ctx, {
		type: 'line',
		data: {
			labels: [],
			datasets: [{
				label: label,
				data: data,
				borderWidth: 1,
				pointRadius: 0,
			}
		]},
		options: {
			scales: {
				y: {
					beginAtZero: true,
					ticks: {
						maxTicksLimit: 8,
					}
				},
				x: {
					grid: {
						display: false
					},
					ticks: {
						callback: () => ('')
					}
				}
			},
			interaction: {
				mode: 'nearest',
				axis: 'x',
				intersect: false,
			},
			legend: {
				labels: {
					fontColor: 'white'
				}
			},
			responsive: true,
			animation: false,
		},
	});
}
