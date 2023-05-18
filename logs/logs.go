package logs

import (
	"time"
	"fmt"
	"os"
	"encoding/json"

	"github.com/rprobaina/lpfs"
)

const logfile = "logs.txt"

type jsonObject struct {
	Date string
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

		msg := jsonObject {
			Date: t,
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
