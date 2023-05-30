const response = await fetch("http://localhost:8080/labels");
const jsonData = await response.json();

let eventSource = new EventSource("/sse");
let currentPid = 1;

addListeners();

let selectRange = document.getElementById('opts');

selectRange.addEventListener('change', function() {
	if (this.value != "runtime") {
		eventSource.close();

		clearPlotData(ptimePlot);
		clearPlotData(memoryPlot);

		fetch(`http://localhost:8080/logs?interval=${this.value}`)
			.then((response) => response.json())
			.then((data) => {
				var map = new Map();

				// Search all pid along log file
				for (let i = 0; i < data.length; ++i) {
					for (let j = 0; j < data[i].Processes.length; ++j) {
						map.set(data[i].Processes[j].Pid, data[i].Processes[j].Comm)
					}
				}

				let rows = document.getElementById("rows")
				rows.innerHTML = "";

				// Create table row for each process
				for (let item of map) {
					const tr = document.createElement('tr');

					const pid = document.createElement('td');
					pid.innerText = item[0];
					tr.append(pid);

					const comm = document.createElement('td');
					comm.innerText = item[1];
					tr.append(comm);

					tr.onclick = function() {
						clearPlotData(ptimePlot);
						clearPlotData(memoryPlot);

						for (let i = 0; i < data.length; ++i) {
							for (let j = 0; j < data[i].Processes.length; ++j) {
								if (data[i].Processes[j].Pid == item[0]) {
									addPlotData(ptimePlot, data[i].Processes[j]);
									break
								}
							}
						}
					}
					rows.appendChild(tr);
				}

			})
			.catch(console.error);
	} else {
		clearPlotData(ptimePlot);
		clearPlotData(memoryPlot);

		eventSource = new EventSource("/sse");
		addListeners();
	}
});

function addListeners() {
	eventSource.addEventListener('memory', e => {
		//console.log(e);
	});

	eventSource.addEventListener('uptime', e => {
		//console.log(e);
	});

	eventSource.addEventListener('swap', e => {
		//console.log(e);
	});

	eventSource.addEventListener('process', e => {
		let obj = JSON.parse(e.data);
		//console.log(obj[currentPid]);

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

			tr.onclick = function() {
				clearPlotData(ptimePlot);
				clearPlotData(memoryPlot);

				currentPid = i;
			}

			rows.appendChild(tr);
		}

		addPlotData(memoryPlot, obj[currentPid].Pid);
		addPlotData(ptimePlot, obj[currentPid]);
	});
}

function clearPlotData(chart) {

	if (chart == ptimePlot) {
		for (let j = 0; j < 4; ++j)
			chart.data.datasets[j].data = [];
		chart.data.labels = [];
	} else {
		chart.data.datasets[0].data = [];
		chart.data.labels = [];
	}

	chart.update();
}

function addPlotData(chart, data) {
	const labels = chart.data.labels

	if (chart == ptimePlot) {
		const array = chart.data.datasets

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
	} else {
		const array = chart.data.datasets[0].data

		if (array.length > 50) {
			array.shift();
			labels.shift();
		}

		array.push(data);
		labels.push("1")
	}

	chart.update();
}
