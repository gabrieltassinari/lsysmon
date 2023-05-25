const response = await fetch("http://localhost:8080/labels");
const jsonData = await response.json();

const eventSource = new EventSource("/sse");

var selectRange= document.getElementById('opts');

selectRange.addEventListener('change', function() {
	if (this.value != "runtime") {
		eventSource.close()

		const jsonData = fetch(`http://localhost:8080/logs?interval=${this.value}`)
			.then((response) => {
				return response.json();
			})

		// TODO: Plot data from request

	} else {
		eventSource = new EventSource("/sse");
	}
});

eventSource.addEventListener('memory', e => {
	//console.log(e);
});

eventSource.addEventListener('uptime', e => {
	//console.log(e);
});

eventSource.addEventListener('swap', e => {
	//console.log(e);
});

let currentPid = 1;

eventSource.addEventListener('process', e => {
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

			for (let j = 0; j < 4; ++j)
				ptimePlot.data.datasets[j].data = [];
			ptimePlot.data.labels = [];

			currentPid = i;
		}

		rows.appendChild(tr);
	}

	// TODO: Append data in plots
	addData(memoryPlot, obj[currentPid].Pid);
	ptimePlotData(obj[currentPid]);
});

function ptimePlotData(data) {
	const labels = ptimePlot.data.labels
	const array = ptimePlot.data.datasets

	if (array[0].data.length > 50) {
		for (let i = 0; i < 4; ++i)
			array[i].data.shift();
		labels.shift();
	}

	array[0].data.push(data.Utime);
	array[1].data.push(data.Stime);
	array[2].data.push(data.Cutime);
	array[3].data.push(data.Cstime);
	labels.push("1");

	ptimePlot.update();
}

function addData(chart, data) {
	const labels = chart.data.labels
	const array = chart.data.datasets[0].data

	if (array.length > 50) {
		labels.shift();
		array.shift();
	}

	array.push(data);
	labels.push("1");

	chart.update();
}
