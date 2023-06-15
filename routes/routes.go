package routes

import (
	"log"
	"net/http"
	"time"

	"lsysmon/data"
	"lsysmon/logs"
)

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		err := data.Memory(w)
		if err != nil {
			log.Println(err)
			return
		}

		err = data.Uptime(w)
		if err != nil {
			log.Println(err)
			return
		}

		err = data.Swap(w)
		if err != nil {
			log.Println(err)
			return
		}

		err = data.ProcessesStat(w)
		if err != nil {
			log.Println(err)
			return
		}

		w.(http.Flusher).Flush()
		time.Sleep(time.Second)
	}
}

func labelsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := data.Labels(w)
	if err != nil {
		log.Println(err)
		return
	}
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	interval := r.URL.Query().Get("interval")
	if interval != "day" && interval != "week" && interval != "month" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := logs.LogsRead(w, interval)
	if err != nil {
		log.Println(err)
		return
	}
}

func Routes() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/sse", sseHandler)
	http.HandleFunc("/labels", labelsHandler)
	http.HandleFunc("/logs", logsHandler)

	errs := make(chan error, 1)
	go logs.Logs(errs)
	go func() {
		for {
			select {
			case err := <-errs:
				log.Println(err)
			}
		}
	}()
}
