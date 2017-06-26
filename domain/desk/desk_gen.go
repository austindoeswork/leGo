// go run gen/gen.go
// VERY GENERATED PLZ NO MODIFY

// Package Desk
// Desk where car keys can be stored
package desk

import (
	"encoding/json"
	"fmt"
	"time"

	"git.ottoq.com/otto-backend/valet/entity"
)

type Desk struct {
	ID        string
	TypeID    string
	Timestamp time.Time
	Name      string
	Lat       float64
	Lng       float64
	NodeID    string
}

func New(
	name string,
	lat float64,
	lng float64,
	nodeID string,
) (*Desk, error) {
	d := &Desk{
		ID:        entity.UUID(),
		TypeID:    "E1874C161CDB492FB95EF210E653B886",
		Timestamp: entity.Now(),
		Name:      name,
		Lat:       lat,
		Lng:       lng,
		NodeID:    nodeID,
	}
	return d, nil
}

// func NewFromRow(row *sql.Row) (*Desk, error) {
// d := Desk
// err := row.Scan(
// //////////////
// }

func Schema() string {
	return `CREATE TABLE Desk (
ID BINARY(16),
TypeID BINARY(16),
Timestamp DATETIME,
Name VARCHAR(100),
Lat FLOAT,
Lng FLOAT,
NodeID BINARY(16),
PRIMARY KEY (ID),
FOREIGN KEY (NodeID) REFERENCES Node(ID)
); `
}

func TableName() string {
	return "Desk"
}

func Random() *Desk {
	d := &Desk{
		ID:        entity.UUID(),
		TypeID:    "E1874C161CDB492FB95EF210E653B886",
		Timestamp: entity.Now(),
		Name:      entity.RANDstring(),
		Lat:       entity.RANDfloat64(),
		Lng:       entity.RANDfloat64(),
		NodeID:    entity.RANDstring(),
	}
	return d
}

func (o *Desk) InsertString() string {
	istr := fmt.Sprintf(`INSERT INTO Desk VALUES(
UNHEX( '%s' ),
UNHEX( '%s' ),
'%s',
'%s',
%f,
%f,
UNHEX( '%s' )
);`,
		o.ID,
		o.TypeID,
		o.Timestamp.Format("2006-01-02 15:04:05"),
		o.Name,
		o.Lat,
		o.Lng,
		o.NodeID,
	)
	return istr
}

func (o *Desk) String() string {
	b, _ := json.MarshalIndent(o, "", "    ")
	return string(b)
}

func (o *Desk) PPrint() {
	fmt.Println(o.String())
}
