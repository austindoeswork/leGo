//go run gen/gen.go
// VERY GENERATED PLZ NO MODIFY

// Package database
// Database persists domain objects
package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func New(addr, dbname, user, pass string) (*Database, error) {
	var err error
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?parseTime=true", user, pass, addr, dbname))
	if err != nil {
		return nil, err
	}
	err = EnsureTablesExist(db)
	if err != nil {
		return nil, err
	}
	return &Database{
		db: db,
	}, nil
}

func EnsureTablesExist(db *sql.DB) error {
	tablequery := "SHOW TABLES LIKE '%s'"
	for _, ts := range Tables {
		var table string
		err := db.QueryRow(fmt.Sprintf(tablequery, ts.Table)).Scan(&table)
		if err != nil {
			_, err := db.Exec(ts.Schema)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("created TABLE %s\n", ts.Table)
		}
	}
	return nil
}

// TableSchema holds an association between a table and its schema
type TableSchema struct {
	Table  string
	Schema string
}

// Tables is an array of TableSchemas
var Tables = []TableSchema{
	TableSchema{
		Table: "Node",
		Schema: `CREATE TABLE Node (
ID BINARY(16),
TypeID BINARY(16),
Timestamp DATETIME,
Name VARCHAR(100),
PRIMARY KEY (ID)
);`,
	},
	TableSchema{
		Table: "Desk",
		Schema: `CREATE TABLE Desk (
ID BINARY(16),
TypeID BINARY(16),
Timestamp DATETIME,
Name VARCHAR(100),
Lat FLOAT,
Lng FLOAT,
NodeID BINARY(16),
PRIMARY KEY (ID),
FOREIGN KEY (NodeID) REFERENCES Node(ID)
);`,
	},
}
