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

		err = data.ProcessesStat(w)
		if err != nil {
			return
		}

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

func logsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return

	}

	/*
	Query URL to get interval value
	f := r.URL.Query()
	interval := f["interval"]
	*/

	// TODO: Return filtered data
	w.Write([]byte("logsData\n"))
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/sse", sse)
	http.HandleFunc("/labels", labels)


	errs := make(chan error, 1)
	go logs.Logs(errs)
	go func() {
		for {
			err := <-errs
			fmt.Println(err)
		}
	}()
	http.HandleFunc("/logs", logsHandler)

	logs.LogsRead("2023-05-24 16:56:34", errs)

	http.ListenAndServe(":8080", nil)
}
