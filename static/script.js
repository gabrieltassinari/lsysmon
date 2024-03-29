let eventSource = new EventSource("/sse");
var currentPid = 1;
var interval = "";

addListeners();

let selectRange = document.getElementById('opts');

selectRange.addEventListener('change', function() {
	interval = this.value;
	if (interval != "runtime") {
		eventSource.close();

		clearPlotData(ptimePlot);
		clearPlotData(memoryPlot);

		fetch(`http://localhost:8080/logs?interval=${interval}`)
			.then((response) => response.json())
			.then((data) => {
				var map = new Map();

				for (let i = 0; i < data.length; ++i)
					map.set(data[i].Pid, data[i].Comm);

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

						fetch(`http://localhost:8080/logs?interval=${interval}&pid=${item[0]}`)
							.then((response) => response.json())
							.then((data) => {
								ptimePlot.data.datasets[0].data = data["Utime"]
								ptimePlot.data.datasets[1].data = data["Stime"]
								ptimePlot.data.datasets[2].data = data["Cutime"]
								ptimePlot.data.datasets[3].data = data["Cstime"]
								ptimePlot.data.labels = data["Date"]
								ptimePlot.update();
							});
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
	eventSource.addEventListener('cpu', e => {
		cpuPlot.data.datasets[0].data.push(e.data);
		cpuPlot.data.labels.push('');
		cpuPlot.update();
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

		let date = new Date();

		let year = dateZero(date.getFullYear()+"");
		let month = dateZero(date.getMonth()+"");
		let day = dateZero(date.getDate()+"");

		let hour = dateZero(date.getHours()+"")
		let minutes = dateZero(date.getMinutes()+"")
		let seconds = dateZero(date.getSeconds()+"")

		let fdate = year + "-" + month + "-" + day  + " "
		let datetime = fdate + hour + ":" + minutes + ":" + seconds

		addPlotData(memoryPlot, obj[currentPid].Pid, datetime);
		addPlotData(ptimePlot, obj[currentPid], datetime);
	});
}

function dateZero(date) {
	if (date.length == 1)
		return "0" + date;
	return date
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

function addPlotData(chart, data, label) {
	const labels = chart.data.labels
	console.log(label)

	if (chart == ptimePlot) {
		const array = chart.data.datasets

		if (array[0].data.length > 50)
			for (let i = 0; i < 4; ++i)
				array[i].data.shift();

		array[0].data.push(data.Utime);
		array[1].data.push(data.Stime);
		array[2].data.push(data.Cutime);
		array[3].data.push(data.Cstime);

	} else {
		const array = chart.data.datasets[0].data

		if (array.length > 50)
			array.shift();
		array.push(data);
	}

	if (labels.length > 50)
		labels.shift();

	labels.push(label);

	chart.update();
}
