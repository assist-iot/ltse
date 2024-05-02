package apisql

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("postgres", getDbURL())

	if err != nil {
		// log.Fatal(err)
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func GetDB() (*sql.DB, error) {
	err := InitDB()
	if err != nil {
		return nil, err
	}

	return db, nil
}
