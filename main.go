package main

import (
    "fmt"
    "time"
    "net/http"

    "lsysmon/data"
)

func sse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {

		memorydata, err := data.Memory()
		if err != nil {
			return
		}
		fmt.Fprintf(w, memorydata)

		uptimedata, err := data.Uptime()
		if err != nil {
			return
		}
		fmt.Fprintf(w, uptimedata)

		swapdata, err := data.Swap()
		if err != nil {
			return
		}

		// Listen data on cli
		fmt.Print(memorydata)
		fmt.Print(uptimedata)
		fmt.Print(swapdata)

		w.(http.Flusher).Flush()
		time.Sleep(time.Second)
	}
}

func main() {
    http.Handle("/", http.FileServer(http.Dir("static")))
    http.HandleFunc("/sse", sse)

    http.ListenAndServe(":8080", nil)
}
