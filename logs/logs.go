package logs

import (
	"encoding/json"
	"fmt"
	"net/http"
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

		db, err := database.OpenConnection()
		if err != nil {
			errs <- fmt.Errorf("Write: %v", err)
			continue
		}

		defer db.Close()

		p, err := lpfs.GetPerProcessStat()
		if err != nil {
			errs <- fmt.Errorf("Write: %v", err)
			continue
		}

		b, err := json.Marshal(p)
		if err != nil {
			errs <- fmt.Errorf("Write: %v", err)
		}

		sql := `INSERT INTO processes (processes_date, processes_stat) VALUES ($1, $2)`

		_, err = db.Exec(sql, time.Now(), string(b))
		if err != nil {
			errs <- fmt.Errorf("Query: %v", err)
			continue
		}
	}
}

func ReadProcess(w http.ResponseWriter, interval string, pid string) error {
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

	b, err := json.Marshal(processes)
	if err != nil {
		return err
	}

	w.Write(b)

	return nil
}
