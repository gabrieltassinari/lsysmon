package routes

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gabrieltassinari/lsysmon/runtime"
	"github.com/rprobaina/lpfs"
)

func TestSse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(sseHandler))

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Errorf("unable to create a request, err: %v", err)
	}

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("failed to request to sse, code: %v, err: %v", res.StatusCode, err)
	}

	if http.StatusOK != res.StatusCode {
		t.Errorf("status: %v, code: %v", http.StatusOK, res.StatusCode)
	}

	buffer := make([]byte, 0, bufio.MaxScanTokenSize)

	scanner := bufio.NewScanner(res.Body)
	scanner.Buffer(buffer, 1024*256)

	for scanner.Scan() {
		if scanner.Text() == "event: memory" {
			scanner.Scan()

			var m runtime.MemoryJSON

			err := json.Unmarshal(scanner.Bytes()[6:], &m)
			if err != nil {
				t.Errorf("invalid response from memory event, err: %v", err)
			}
		}

		if scanner.Text() == "event: uptime" {
			scanner.Scan()
			u := scanner.Text()[6:]

			_, err := strconv.ParseFloat(u, 64)
			if err != nil {
				t.Errorf("invalid response from uptime event, err: %v", err)
			}
		}

		if scanner.Text() == "event: swap" {
			scanner.Scan()

			var s runtime.SwapJSON

			err := json.Unmarshal(scanner.Bytes()[6:], &s)
			if err != nil {
				t.Errorf("invalid response from swap event, err: %v", err)
			}
		}

		if scanner.Text() == "event: process" {
			scanner.Scan()

			var p []lpfs.Procstat

			err := json.Unmarshal(scanner.Bytes()[6:], &p)
			if err != nil {
				t.Errorf("invalid response from process event, err: %v", err)
			}

			break
		}

	}

	if err := scanner.Err(); err != nil {
		t.Errorf("unable to read events, err: %v", err)
	}

	fmt.Printf("sseHandler(): %v\n", res.StatusCode)
}

func TestReadProcess(t *testing.T) {
	rrd := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/logs?interval=day&pid=1", nil)
	logsHandler(rrd, req)

	if http.StatusOK != rrd.Code {
		t.Errorf("status: %v, code: %v", http.StatusOK, rrd.Code)
	}
	fmt.Printf("ReadProcess() day: %v\n", rrd.Code)

	rrw := httptest.NewRecorder()

	req = httptest.NewRequest(http.MethodGet, "/logs?interval=week&pid=1", nil)
	logsHandler(rrw, req)

	if http.StatusOK != rrw.Code {
		t.Errorf("status: %v, code: %v", http.StatusOK, rrw.Code)
	}
	fmt.Printf("ReadProcess() week: %v\n", rrw.Code)

	rrm := httptest.NewRecorder()

	req = httptest.NewRequest(http.MethodGet, "/logs?interval=month&pid=1", nil)
	logsHandler(rrm, req)

	if http.StatusOK != rrm.Code {
		t.Errorf("status: %v, code: %v", http.StatusOK, rrm.Code)
	}
	fmt.Printf("ReadProcess() month: %v\n", rrm.Code)

	rri := httptest.NewRecorder()

	req = httptest.NewRequest(http.MethodGet, "/logs?interval=month&pid=error", nil)
	logsHandler(rri, req)

	if http.StatusOK == rri.Code {
		t.Errorf("status: %v, code: %v", http.StatusBadRequest, rri.Code)
	}
	fmt.Printf("ReadProcess() invalid: %v\n", rri.Code)

	rrj := httptest.NewRecorder()

	req = httptest.NewRequest(http.MethodGet, "/logs?interval=month&pid=4194305", nil)
	logsHandler(rrj, req)

	if http.StatusOK == rri.Code {
		t.Errorf("status: %v, code: %v", http.StatusBadRequest, rrj.Code)
	}
	fmt.Printf("ReadProcess() invalid: %v\n", rrj.Code)

}

func TestReadProcesses(t *testing.T) {
	rrd := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/logs?interval=day", nil)
	logsHandler(rrd, req)

	if http.StatusOK != rrd.Code {
		t.Errorf("status: %v, code: %v", http.StatusOK, rrd.Code)
	}
	fmt.Printf("ReadProcesses() day: %v\n", rrd.Code)

	rrw := httptest.NewRecorder()

	req = httptest.NewRequest(http.MethodGet, "/logs?interval=week", nil)
	logsHandler(rrw, req)

	if http.StatusOK != rrw.Code {
		t.Errorf("status: %v, code: %v", http.StatusOK, rrw.Code)
	}
	fmt.Printf("ReadProcesses() week: %v\n", rrw.Code)

	rrm := httptest.NewRecorder()

	req = httptest.NewRequest(http.MethodGet, "/logs?interval=month", nil)
	logsHandler(rrm, req)

	if http.StatusOK != rrm.Code {
		t.Errorf("status: %v, code: %v", http.StatusOK, rrm.Code)
	}
	fmt.Printf("ReadProcesses() month: %v\n", rrm.Code)

	rri := httptest.NewRecorder()

	req = httptest.NewRequest(http.MethodGet, "/logs?interval=invalid", nil)
	logsHandler(rri, req)

	if http.StatusOK == rri.Code {
		t.Errorf("status: %v, code: %v", http.StatusBadRequest, rri.Code)
	}
	fmt.Printf("ReadProcesses() invalid: %v\n", rri.Code)
}
