const ctx = document.getElementById('myChart');
const plot = new Chart(ctx, {
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
			}
		},
		responsive: true,
	},
});

const eventSource = new EventSource("/sse");

eventSource.addEventListener('memory', e => {
	let json = e.data;
	let obj = JSON.parse(json);

	addData(plot, obj);

	document.getElementById('sse').innerText = (obj.Free/1048576).toFixed(2);
});

function addData(chart, data) {
	const labels = chart.data.labels
	const array = chart.data.datasets[0].data

	if (array.length > 50) {
		labels.shift();
		array.shift();
	}

	labels.push("1");
	array.push((data.Free/1048576).toFixed(2));

	chart.update();
}
