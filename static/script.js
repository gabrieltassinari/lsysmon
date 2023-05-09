const response = await fetch("http://localhost:8080/labels");
const jsonData = await response.json();
console.log(jsonData);

const eventSource = new EventSource("/sse");
/*
eventSource.addEventListener('memory', e => {
	console.log(e);

	let json = e.data;
	let obj = JSON.parse(json);

	addData(memoryPlot, (obj.Free/1048576).toFixed(2));
});
*/
eventSource.addEventListener('uptime', e => {
	console.log(e);
});

eventSource.addEventListener('swap', e => {
	console.log(e);

	let json = e.data;
	let obj = JSON.parse(json);

	addData(swapPlot, obj.Used/1048576);
});


let procMemory = 1;

eventSource.addEventListener('process', e => {
	console.log(e);
	let obj = JSON.parse(e.data);

	let rows = document.getElementById("rows")

	rows.innerHTML = "";

	// Create a table with all processes
	for (let i = 0; i < obj.length; ++i) {
		const tr = document.createElement('tr');

		const pid = document.createElement('td');
		pid.innerText = `${obj[i].Pid}`;
		tr.append(pid);

		const comm = document.createElement('td');
		comm.innerText = `${obj[i].Comm}`;
		tr.appendChild(comm);

		const state = document.createElement('td');
		state.innerText = `${obj[i].State}`;
		tr.appendChild(state);

		// TODO: Create plots with process stat information
		tr.onclick = function() {
			memoryPlot.data.datasets[0].data = [];
			memoryPlot.data.labels = [];
			console.log(`${obj[i].Pid}`);

			procMemory = obj[i].Pid;
		}

		rows.appendChild(tr);
	}

	// TODO: Append data in plots
	addData(memoryPlot, obj[procMemory].Utime);
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
