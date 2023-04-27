const ctx = document.getElementById('myChart');

const config = {
	type: 'line',
	data: {
		labels: [],
		datasets: [{
			label: 'Free memory in gigabytes',
			data: [],
			fill: true,
			borderWidth: 1
		}
	]},
	options : {
		scales: {
			y: {
				suggestedMin: 0,
				suggestedMax: 16
			},
			x: {
				ticks: {
					display: false
				}
			}
		},
		responsive: true,
	},
};

const plot = new Chart(ctx, config)

const eventSource = new EventSource("/sse");

eventSource.addEventListener('memory', e => {
	console.log(e);

	let json = e.data;
	let obj = JSON.parse(json);

	addData(plot, (obj.Free/1048576).toFixed(2));

	document.getElementById('sse').innerText = (obj.Free/1048576).toFixed(2);
});

eventSource.addEventListener('uptime', e => {
	console.log(e);
});

function addData(chart, data) {
	const labels = chart.data.labels
	const array = chart.data.datasets[0].data

	if (array.length > 50) {
		labels.shift();
		array.shift();
	}

	labels.push("1");
	array.push(data);

	chart.update();
}
