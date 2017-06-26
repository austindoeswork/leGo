package main

import (
	"fmt"
	"log"
	"time"

	"git.ottoq.com/otto-backend/valet/database"
	"git.ottoq.com/otto-backend/valet/domain/desk"
	"git.ottoq.com/otto-backend/valet/domain/node"
)

func main() {
	var err error
	dbadd := "tcp(127.0.0.1:3306)"
	db, err := database.New(dbadd, "test3", "austin", "")
	if err != nil {
		log.Fatal(err)
	}

	n := node.Random()
	d := desk.Random()
	d.NodeID = n.ID
	// n.PPrint()
	// d.PPrint()

	_, err = db.Exec(n.InsertString())
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(d.InsertString())
	if err != nil {
		log.Fatal(err)
	}

	// rows, err := db.Query("SELECT HEX(ID) as ID from Node")
	rows, err := db.Query("SELECT * from Node")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var id string
		var typeid string
		var name string
		var timestamp time.Time
		// err = rows.Scan(&id)
		err = rows.Scan(&id, &typeid, &timestamp, &name)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println(id)
		fmt.Println(id, typeid, name)
	}
	// d.PPrint()
	// var id, typeid string
}
