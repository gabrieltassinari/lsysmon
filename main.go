package main

import (
    "fmt"
    "time"
    "net/http"
    "encoding/json"

    "github.com/rprobaina/lpfs"
)

type Memory struct{
	Total int //`json:"Total"`
	Free int //`json:"Free"`
}

func sse(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    for {
	mt, _ := lpfs.GetMemTotal()
	mf, _ := lpfs.GetMemFree()

	msg := Memory{
		Total: mt, 
		Free: mf, 
	}

	b, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	
	fmt.Fprintf(w, "event: memory\n")
        fmt.Fprintf(w, "data: %v\n\n", string(b))
        w.(http.Flusher).Flush()

	time.Sleep(time.Second)
    }
}

func main() {
    http.Handle("/", http.FileServer(http.Dir("static")))
    http.HandleFunc("/sse", sse)

    http.ListenAndServe(":8080", nil)
}
