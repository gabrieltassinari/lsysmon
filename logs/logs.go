package logs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gabrieltassinari/lsysmon/database"
	"github.com/rprobaina/lpfs"
)

func readInterval(t time.Time, interval string) (string, error) {
	var start string

	switch interval {
	case "day":
		start = t.AddDate(0, 0, -1).Format(time.DateTime)
	case "week":
		start = t.AddDate(0, 0, -7).Format(time.DateTime)
	case "month":
		start = t.AddDate(0, -1, 0).Format(time.DateTime)
	default:
		return "", fmt.Errorf("Invalid time interval.")
	}

	return start, nil
}

func WriteLogs(errs chan error) {
	for {
		time.Sleep(time.Minute)

		// TODO: Create a directory to store logs.
		file := fmt.Sprintf("./logs/%v.log", time.Now().Format(time.DateOnly))

		f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			errs <- fmt.Errorf("WriteLogs: %v", err)
			continue
		}

		defer f.Close()

		p, err := lpfs.GetPerProcessStat()
		if err != nil {
			errs <- fmt.Errorf("WriteLogs: %v", err)
			continue
		}

		b, err := json.Marshal(p)
		if err != nil {
			errs <- fmt.Errorf("WriteLogs: %v", err)
			continue
		}

		// Write func doesn't have breakline by default
		b = append(b, 10)

		f.Write(b)
		if err != nil {
			f.Close()
			errs <- fmt.Errorf("WriteLogs: %v", err)
			continue
		}

		err = f.Close()
		if err != nil {
			errs <- fmt.Errorf("WriteLogs: %v", err)
			continue
		}

		fmt.Printf("WriteLogs: Sucess writing to %v\n", file)
	}
}

func WriteCpuUsage(errs chan error) {
	for {
		time.Sleep(time.Minute)

		var prevIdle, prevTotal int
		var cpuUsage float64

		db, err := database.OpenConnection()
		if err != nil {
			errs <- fmt.Errorf("Write cpu: %v", err)
			continue
		}

		for i := 0; i < 2; i++ {
			user, err := lpfs.GetCpuUserTime()
			if err != nil {
				errs <- fmt.Errorf("Write cpu: %v", err)
				continue
			}

			nice, err := lpfs.GetCpuNiceTime()
			if err != nil {
				errs <- fmt.Errorf("Write cpu: %v", err)
				continue
			}

			system, err := lpfs.GetCpuSystemTime()
			if err != nil {
				errs <- fmt.Errorf("Write cpu: %v", err)
				continue
			}

			idle, err := lpfs.GetCpuIdleTime()
			if err != nil {
				errs <- fmt.Errorf("Write cpu: %v", err)
				continue
			}

			iowait, err := lpfs.GetCpuIowaitTime()
			if err != nil {
				errs <- fmt.Errorf("Write cpu: %v", err)
				continue
			}

			irq, err := lpfs.GetCpuIrqTime()
			if err != nil {
				errs <- fmt.Errorf("Write cpu: %v", err)
				continue
			}

			softirq, err := lpfs.GetCpuSoftirqTime()
			if err != nil {
				errs <- fmt.Errorf("Write cpu: %v", err)
				continue
			}

			steal, err := lpfs.GetCpuStealTime()
			if err != nil {
				errs <- fmt.Errorf("Write cpu: %v", err)
				continue
			}

			total := user + nice + system + idle + iowait + irq + softirq + steal

			if i > 0 {
				deltaIdle := idle - prevIdle
				deltaTotal := total - prevTotal
				cpuUsage = (1.0 - float64(deltaIdle)/float64(deltaTotal)) * 100.0
			}

			prevIdle = idle
			prevTotal = total
			time.Sleep(time.Second)
		}

		sql := `INSERT INTO cpu (cpu_date, cpu_usage) VALUES ($1, $2)`

		_, err = db.Exec(sql, time.Now(), cpuUsage)
		if err != nil {
			errs <- fmt.Errorf("Insert cpu: %v", err)
		}
	}
}

func ReadProcess(w http.ResponseWriter, interval string, pids string) error {
	db, err := database.OpenConnection()
	if err != nil {
		return err
	}

	end := time.Now()

	start, err := readInterval(end, interval)
	if err != nil {
		return err
	}

	pid, err := strconv.Atoi(pids)
	if err != nil {
		return err
	}

	sql :=
		`SELECT
			(elem->'Utime'),
			(elem->'Stime'),
			(elem->'Cutime'),
			(elem->'Cstime'),
			processes_date
		 FROM
			processes
		 CROSS JOIN
			jsonb_array_elements(processes_stat) elem
		 WHERE
			(elem->>'Pid') = $1
			AND processes_date > $2
			AND processes_date < $3`

	rows, err := db.Query(sql, pid, start, end.Format(time.DateTime))
	if err != nil {
		return err
	}
	defer rows.Close()

	type Process struct {
		Utime  []int
		Stime  []int
		Cutime []int
		Cstime []int
		Date   []string
	}

	var p Process

	for rows.Next() {
		var (
			utime  int
			stime  int
			cutime int
			cstime int
			date   time.Time
		)

		if err := rows.Scan(&utime, &stime, &cutime, &cstime, &date); err != nil {
			return err
		}

		p.Utime = append(p.Utime, utime)
		p.Stime = append(p.Utime, stime)
		p.Cutime = append(p.Utime, cutime)
		p.Cstime = append(p.Utime, cstime)
		p.Date = append(p.Date, date.Format(time.DateTime))

	}

	if len(p.Utime) == 0 {
		return fmt.Errorf("No results were returned by the query")
	}

	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}

func ReadProcesses(w http.ResponseWriter, interval string) error {
	db, err := database.OpenConnection()
	if err != nil {
		return err
	}

	end := time.Now()

	start, err := readInterval(end, interval)
	if err != nil {
		return err
	}

	sql :=
		`SELECT
			(elem->>'Pid'),
			(elem->>'Comm'),
			(elem->>'State')
		 FROM
			processes
		 CROSS JOIN
			jsonb_array_elements(processes_stat) elem
		 WHERE
			processes_date >= $1
			AND processes_date < $2`

	rows, err := db.Query(sql, start, end.Format(time.DateTime))
	if err != nil {
		return err
	}
	defer rows.Close()

	type Process struct {
		Pid   string
		Comm  string
		State string
	}

	var processes []Process

	for rows.Next() {
		var scans Process

		if err := rows.Scan(&scans.Pid, &scans.Comm, &scans.State); err != nil {
			return err
		}

		processes = append(processes, scans)
	}

	if len(processes) == 0 {
		return fmt.Errorf("No results were returned by the query")
	}

	b, err := json.Marshal(processes)
	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}

func ReadCpuUsage(w http.ResponseWriter, interval string) error {
	db, err := database.OpenConnection()
	if err != nil {
		return err
	}

	end := time.Now()

	start, err := readInterval(end, interval)
	if err != nil {
		return err
	}

	sql :=
		`SELECT
			cpu_date,
			cpu_usage
		 FROM
			cpu
		 WHERE
			cpu_date >= $1 AND
			cpu_date < $2`

	rows, err := db.Query(sql, start, end.Format(time.DateTime))

	type Cpu struct {
		Date  []string
		Usage []float64
	}

	var c Cpu

	for rows.Next() {
		var (
			date  time.Time
			usage float64
		)

		if err := rows.Scan(&date, &usage); err != nil {
			return err
		}

		c.Date = append(c.Date, date.Format(time.DateTime))
		c.Usage = append(c.Usage, usage)
	}

	if len(c.Usage) == 0 {
		return fmt.Errorf("No results were returned by the query")
	}

	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}
