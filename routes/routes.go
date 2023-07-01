package routes

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gabrieltassinari/lsysmon/logs"
	"github.com/gabrieltassinari/lsysmon/runtime"
)

func sseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		err := runtime.Memory(w)
		if err != nil {
			log.Println(err)
			return
		}

		err = runtime.Uptime(w)
		if err != nil {
			log.Println(err)
			return
		}

		err = runtime.Swap(w)
		if err != nil {
			log.Println(err)
			return
		}

		err = runtime.ProcessesStat(w)
		if err != nil {
			log.Println(err)
			return
		}

		w.(http.Flusher).Flush()
		time.Sleep(time.Second)
	}
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	interval := r.URL.Query().Get("interval")
	fmt.Println("Interval: ", interval)
	if interval != "day" && interval != "week" && interval != "month" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err := logs.ReadProcesses(w, interval)
	if err != nil {
		log.Println(err)
		return
	}
}

func Routes() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/sse", sseHandler)
	http.HandleFunc("/logs", logsHandler)

	errs := make(chan error, 1)
	go logs.WriteLogs(errs)
	go func() {
		for {
			select {
			case err := <-errs:
				log.Println(err)
			}
		}
	}()
}
