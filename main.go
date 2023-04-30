package main

import (
    "time"
    "net/http"

    "lsysmon/data"
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

		w.(http.Flusher).Flush()
		time.Sleep(time.Second)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/sse", sse)

	http.ListenAndServe(":8080", nil)
}
