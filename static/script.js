const response = await fetch("http://localhost:8080/labels");
const jsonData = await response.json();
console.log(jsonData);

const eventSource = new EventSource("/sse");

eventSource.addEventListener('memory', e => {
	console.log(e);

	let json = e.data;
	let obj = JSON.parse(json);

	addData(memoryPlot, (obj.Free/1048576).toFixed(2));
});

eventSource.addEventListener('uptime', e => {
	console.log(e);
});

eventSource.addEventListener('swap', e => {
	console.log(e);

	let json = e.data;
	let obj = JSON.parse(json);

	addData(swapPlot, obj.Used/1048576);
});

eventSource.addEventListener('process', e => {
	console.log(e);
	let obj = JSON.parse(e.data);
	console.log(obj[0]);
	console.log(obj[1]);
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
