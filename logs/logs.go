package logs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gabrieltassinari/lsysmon/database"
	"github.com/rprobaina/lpfs"
)

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

		_, err = db.Exec(sql, time.Now(), b)
		if err != nil {
			errs <- fmt.Errorf("Query: %v", err)
			continue
		}
	}
}

func ReadProcesses(w http.ResponseWriter, interval string) error {
	db, err := database.OpenConnection()
	if err != nil {
		return err
	}

	var start string
	end := time.Now()

	switch interval {
	case "day":
		start = end.AddDate(0, 0, -1).Format(time.DateTime)
	case "week":
		start = end.AddDate(0, 0, -7).Format(time.DateTime)
	case "month":
		start = end.AddDate(0, -1, 0).Format(time.DateTime)
	default:
		return err
	}

	sql := `SELECT processes_stat FROM processes WHERE processes_date >= $1 AND processes_date < $2`

	rows, err := db.Query(sql, start, end.Format(time.DateTime))
	if err != nil {
		return err
	}

	var s []byte
	s = append(s, []byte("[")...)

	for rows.Next() {
		var data []byte

		if err := rows.Scan(&data); err != nil {
			return err
		}

		s = append(s, data...)
		s = append(s, []byte(",")...)
	}

	s = s[:len(s)-1]
	s = append(s, []byte("]")...)

	w.Write(s)
	return nil
}
