package database

import (
	"database/sql"
)

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	if len(args) == 0 {
		return d.db.Exec(query)
	}
	return d.db.Exec(query, args)
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if len(args) == 0 {
		return d.db.Query(query)
	}
	return d.db.Query(query, args)
}
