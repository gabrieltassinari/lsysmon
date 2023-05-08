const memoryConfig = {
	type: 'line',
	data: {
		labels: [],
		datasets: [{
			label: '',
			data: [],
			fill: true,
			borderWidth: 1,
			pointRadius: 0,
		}
	]},
	options : {
		scales: {
			y: {
				suggestedMin: 0,
				suggestedMax: 16,
			},
			x: {
				grid: {
					display: false,
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
};

const swapConfig = {
	type: 'line',
	data: {
		labels: [],
		datasets: [{
			label: '',
			data: [],
			fill: true,
			borderWidth: 1,
			pointRadius: 0,
		}
	]},
	options : {
		scales: {
			y: {
				suggestedMin: 0,
				//suggestedMax: 16
			},
			x: {
				grid: {
					display: false,
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
	},
};

const memoryCtx = document.getElementById('memory');
const memoryPlot = new Chart(memoryCtx, memoryConfig)

memoryPlot.config.data.datasets[0].label = 'Free memory (GB/s)';

const swapCtx = document.getElementById('swap');
const swapPlot = new Chart(swapCtx, swapConfig)

swapPlot.config.data.datasets[0].label = 'Swap in use (GB/s)';
