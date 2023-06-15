package logs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rprobaina/lpfs"
)

const logfile = "logs.txt"

type jsonObject struct {
	Date      string
	Processes []lpfs.Procstat
}

func Logs(errs chan error) {
	for {
		time.Sleep(time.Minute)

		file, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			errs <- fmt.Errorf("Logs: unable to open %s file: %v", logfile, err)
			continue
		}

		t := time.Now().Format(time.DateTime)

		p, err := lpfs.GetPerProcessStat()
		if err != nil {
			errs <- fmt.Errorf("Logs: unable to get processes stats: %v", err)
			continue
		}

		msg := jsonObject{
			Date:      t,
			Processes: p,
		}
		b, err := json.Marshal(msg)
		if err != nil {
			errs <- fmt.Errorf("Logs: unable to marshal log data: %v", err)
			continue
		}

		_, err = file.Write(b)
		if err != nil {
			errs <- fmt.Errorf("Logs: unable to write in %s file: %v", logfile, err)
			continue
		}

		_, err = file.WriteString("\n")
		if err != nil {
			errs <- fmt.Errorf("Logs: unable to write in %s file: %v", logfile, err)
			continue
		}

		file.Close()
	}
}

func LogsRead(w http.ResponseWriter, interval string) error {
	file, err := os.Open(logfile)
	if err != nil {
		return fmt.Errorf("unable to open %s file: %v", logfile, err)
	}

	fscanner := bufio.NewScanner(file)

	buffer := make([]byte, 0, bufio.MaxScanTokenSize)
	fscanner.Buffer(buffer, 1024*1024)

	var start time.Time
	var find []byte

	endstr := time.Now().Format(time.DateTime)
	end, err := time.Parse(time.DateTime, endstr)
	if err != nil {
		return fmt.Errorf("unable to parse date in %s file: %v", logfile, err)
	}

	// TODO: Handle when last day/month doesnt have data
	if interval == "day" {
		start = end.AddDate(0, 0, -1)
		find = []byte(start.Format(time.DateOnly)[:10])
	}

	if interval == "week" {
		// TODO: FIX WEEK
		start = end.AddDate(0, 0, -7)
	}

	if interval == "month" {
		start = end.AddDate(0, 0, 0)
		find = []byte(start.Format(time.DateOnly)[:7])
	}

	s := []byte("[")

	// Iterating each line of the file
	for fscanner.Scan() {
		if bytes.Contains(fscanner.Bytes(), find) {
			s = append(s, fscanner.Bytes()...)
			s = append(s, []byte(",")...)
		}
	}

	s = s[:len(s)-1]
	s = append(s, []byte("]")...)

	w.Write([]byte(s))

	return nil
}
