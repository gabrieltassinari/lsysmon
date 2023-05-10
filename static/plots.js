const memoryCtx = document.getElementById('memory');
const memoryPlot = createChart(memoryCtx, [], 'Free memory (GB/s)');

const swapCtx = document.getElementById('swap');
const swapPlot = createChart(swapCtx, [], 'Swap in use (GB/s)');

const utimeCtx = document.getElementById('utime');
const utimePlot = createChart(utimeCtx, [], 'Utime');

const cstimeCtx = document.getElementById('cstime');
const cstimePlot = createChart(cstimeCtx, [], 'cstime');

function createChart(ctx, data, label) {
	return new Chart(ctx, {
		type: 'line',
		data: {
			labels: [],
			datasets: [{
				label: label,
				data: data,
				fill: true,
				borderWidth: 1,
				pointRadius: 0,
			}
		]},
		options: {
			scales: {
				y: {
					suggestedMin: 0,
					suggestedMax: 16,
				},
				x: {
					grid: {
						display: false
					},
					ticks: {
						display: false
					}
				}
			},
			interaction: {
				mode: 'nearest',
				axis: 'x',
				intersect: false,
			},
			responsive: true,
			animation: false,
		},
	});
}
