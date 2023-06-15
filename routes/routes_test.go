package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"strconv"
	"bufio"
	"fmt"

	"lsysmon/data"
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

			var m data.MemoryJSON

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

			var s data.SwapJSON

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

func TestLabels(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/labels", nil)

	labelsHandler(rr, req)

	if http.StatusOK != rr.Code {
		t.Errorf("status: %v, code: %v", http.StatusOK, rr.Code)
	}
	fmt.Printf("labelsHandler(): %v\n", rr.Code)
}

func TestLogs(t *testing.T) {
	rr := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/logs?interval=day", nil)
	logsHandler(rr, req)

	if http.StatusOK != rr.Code {
		t.Errorf("status: %v, code: %v", http.StatusOK, rr.Code)
	}
	fmt.Printf("logsHandler() day: %v\n", rr.Code)

	req = httptest.NewRequest(http.MethodGet, "/logs?interval=week", nil)
	logsHandler(rr, req)

	if http.StatusOK != rr.Code {
		t.Errorf("status: %v, code: %v", http.StatusOK, rr.Code)
	}
	fmt.Printf("logsHandler() week: %v\n", rr.Code)

	req = httptest.NewRequest(http.MethodGet, "/logs?interval=month", nil)
	logsHandler(rr, req)

	if http.StatusOK != rr.Code {
		t.Errorf("status: %v, code: %v", http.StatusOK, rr.Code)
	}
	fmt.Printf("logsHandler() month: %v\n", rr.Code)
	
	req = httptest.NewRequest(http.MethodGet, "/logs?interval=invalid", nil)
	logsHandler(rr, req)

	if http.StatusBadRequest != rr.Code {
		t.Errorf("status: %v, code: %v", http.StatusBadRequest, rr.Code)
	}
	fmt.Printf("logsHandler() invalid: %v\n", rr.Code)
}
