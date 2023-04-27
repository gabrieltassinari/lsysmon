package main
import (
    "fmt"
    "time"
    "net/http"
    "encoding/json"

    "github.com/rprobaina/lpfs"
)

var events string

type Memory struct{
	Total int //`json:"Total"`
	Free int //`json:"Free"`
}

func memoryData() {
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

	event := "memory"
	data := fmt.Sprintf("%v", string(b))

	addEvent(event, data)
}

func uptimeData() {
	uptime, _ := lpfs.GetUptimeSystem()

	event := "uptime"
	data := fmt.Sprintf("%f", uptime)

	addEvent(event, data)
}

func addEvent(event string, data string) {
	e := "event: " + event + "\n"
	d := "data: " + data + "\n\n"

	events += e
	events += d
}

func sse(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")


    for {

	memoryData()
	uptimeData()

	fmt.Printf(events)
	fmt.Fprintf(w, events)

        w.(http.Flusher).Flush()

	time.Sleep(time.Second)
	events = ""
    }
}

func main() {
    http.Handle("/", http.FileServer(http.Dir("static")))
    http.HandleFunc("/sse", sse)

    http.ListenAndServe(":8080", nil)
}
