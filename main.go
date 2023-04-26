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

func memoryData() (string, string) {
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

	event := "event: memory\n"
	data := fmt.Sprintf("data: %v\n\n", string(b))

	return event, data
}

func uptimeData() (string, string) {
	uptime, _ := lpfs.GetUptimeSystem()

	event := "event: uptime\n"
	data := fmt.Sprintf("data: %f\n\n", uptime)

	return event, data
}


func sse(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    for {

	event, data := memoryData()

	fmt.Fprintf(w, event)
	fmt.Fprintf(w, data)

	event, data = uptimeData()

	fmt.Fprintf(w, event)
	fmt.Fprintf(w, data)

	//fmt.Fprintf(w, "event: %v\n", event)
        //fmt.Fprintf(w, "data: %v\n\n", string(b))

        w.(http.Flusher).Flush()

	time.Sleep(time.Second)
    }
}

func main() {
    http.Handle("/", http.FileServer(http.Dir("static")))
    http.HandleFunc("/sse", sse)

    http.ListenAndServe(":8080", nil)
}
