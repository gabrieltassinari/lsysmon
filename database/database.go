package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	user     string = os.Getenv("PGUSER")
	password string = os.Getenv("PGPASSWORD")
	host     string = os.Getenv("PGHOST")
	port     string = os.Getenv("PGPORT")
	dbname   string = os.Getenv("PGDATABASE")
)

func OpenConnection() (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return db, nil
}
