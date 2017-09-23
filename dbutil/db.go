package dbutil

import "database/sql"
import _ "github.com/go-sql-driver/mysql"

var db *sql.DB

func Open() (*sql.DB, error) {
	var err error
	if db == nil {
		db, err = sql.Open("mysql", "root:root@/gcai?charset=utf8")
		return db, err
	}
	return db, nil
}