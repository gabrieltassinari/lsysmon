const eventSource = new EventSource("/sse");
eventSource.addEventListener('message', e => {
	let json = e.data;
	let obj = JSON.parse(json);
	console.log(obj);

	document.getElementById('sse').innerText = (obj.Free/1024).toFixed(2);
});
