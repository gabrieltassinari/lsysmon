package logs

import (
	"time"
	"fmt"
	"os"
	"encoding/json"
	"bufio"

	"github.com/rprobaina/lpfs"
)

const logfile = "logs.txt"

type jsonObject struct {
	Date      string
	Processes []lpfs.Procstat
}

func Logs(errs chan error) {
	for {
		file, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			errs <- fmt.Errorf("Logs: unable to open %s file: %v", logfile, err)
		}

		t := time.Now().Format(time.DateTime)

		p, err := lpfs.GetPerProcessStat()
		if err != nil {
			errs <- fmt.Errorf("Logs: unable to get processes stats: %v", err)
		}

		msg := jsonObject{
			Date:      t,
			Processes: p,
		}
		b, _ := json.Marshal(msg)

		_, err = file.Write(b)
		if err != nil {
			errs <- fmt.Errorf("Logs: unable to write in %s file: %v", logfile, err)
		}

		_, err = file.WriteString("\n")
		if err != nil {
			errs <- fmt.Errorf("Logs: unable to write in %s file: %v", logfile, err)
		}

		file.Close()

		time.Sleep(time.Minute)
	}
}

func LogsRead(interval string, errs chan error) {
	file, err := os.Open(logfile)
	if err != nil {
		errs <- fmt.Errorf("LogsRead: unable to open %s file: %v", logfile, err)
	}

	fscanner := bufio.NewScanner(file)

	buffer := make([]byte, 0, bufio.MaxScanTokenSize)
	fscanner.Buffer(buffer, 1024*1024)

	start, _ := time.Parse(time.DateTime, interval)

	endstr := time.Now().Format(time.DateTime)
	end, _ := time.Parse(time.DateTime, endstr)

	// Iterating each line of the file
	for fscanner.Scan() {
		line := fscanner.Text()
		date, err := time.Parse(time.DateTime, line[9:28])
		if err != nil {
			errs <- fmt.Errorf("LogsRead: unable to parse date from %s file: %v", logfile, err)

		if date.After(start) && date.Before(end) {
			// TODO: Append the line to a string slice
		}
	}
}
