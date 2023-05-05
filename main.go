package main

import (
    "net/http"
    "time"
    "fmt"

    "lsysmon/data"
    "lsysmon/logs"
)

func sse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {

		err := data.Memory(w)
		if err != nil {
			return
		}

		err = data.Uptime(w)
		if err != nil {
			return
		}

		err = data.Swap(w)
		if err != nil {
			return
		}

		logs.Logs()

		w.(http.Flusher).Flush()
		time.Sleep(time.Second)
	}
}

func labels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := data.Labels(w)
	if err != nil {
		s := fmt.Sprintf("%v", err)
		http.Error(w, s, 500)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/sse", sse)
	http.HandleFunc("/labels", labels)

	http.ListenAndServe(":8080", nil)
}
