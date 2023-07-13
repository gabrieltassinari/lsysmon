package routes

import (
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

		if err = runtime.Cpu(w); err != nil {
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
	} else if r.Method == "GET" {
		interval := r.FormValue("interval")
		pid := r.FormValue("pid")

		if pid != "" {
			err := logs.ReadProcess(w, interval, pid)
			if err != nil {
				http.Error(w, "400 - Bad Request!", http.StatusBadRequest)
				log.Println(err)
				return
			}
		} else {
			err := logs.ReadProcesses(w, interval)
			if err != nil {
				http.Error(w, "400 - Bad Request!", http.StatusBadRequest)
				log.Println(err)
				return
			}
		}
	}
}

func Routes() {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/sse", sseHandler)
	http.HandleFunc("/logs", logsHandler)

	errs := make(chan error, 1)

	go logs.WriteLogs(errs)
	go logs.WriteCpuUsage(errs)

	go func() {
		for {
			select {
			case err := <-errs:
				log.Println(err)
			}
		}
	}()
}
