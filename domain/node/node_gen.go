// go run gen/gen.go
// VERY GENERATED PLZ NO MODIFY

// Package Node
// Node represents a node in the organization permission heirarchy tree
package node

import (
	"encoding/json"
	"fmt"
	"time"

	"git.ottoq.com/otto-backend/valet/entity"
)

type Node struct {
	ID        string
	TypeID    string
	Timestamp time.Time
	Name      string
}

func New(
	name string,
) (*Node, error) {
	d := &Node{
		ID:        entity.UUID(),
		TypeID:    "0C74DFC158C646C280BCB0DAF9E015D1",
		Timestamp: entity.Now(),
		Name:      name,
	}
	return d, nil
}

type Scannable interface {
	Scan(dest ...interface{}) error
}

func NewFromRow(row Scannable) (*Node, error) {
	d := Node{}
	err := row.Scan(
		&d.ID,
		&d.TypeID,
		&d.Timestamp,
		&d.Name,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func Schema() string {
	return `CREATE TABLE Node (
ID BINARY(16),
TypeID BINARY(16),
Timestamp DATETIME,
Name VARCHAR(100),
PRIMARY KEY (ID)
); `
}

func TableName() string {
	return "Node"
}

func Random() *Node {
	d := &Node{
		ID:        entity.UUID(),
		TypeID:    "0C74DFC158C646C280BCB0DAF9E015D1",
		Timestamp: entity.Now(),
		Name:      entity.RANDstring(),
	}
	return d
}

func (o *Node) InsertString() string {
	istr := fmt.Sprintf(`INSERT INTO Node VALUES(
UNHEX( '%s' ),
UNHEX( '%s' ),
'%s',
'%s'
)
ON DUPLICATE KEY UPDATE
TypeID=UNHEX( '%s' ),
Timestamp='%s',
Name='%s'
;`,
		o.ID,
		o.TypeID,
		o.Timestamp.Format("2006-01-02 15:04:05"),
		o.Name,

		o.TypeID,
		o.Timestamp.Format("2006-01-02 15:04:05"),
		o.Name,
	)
	return istr
}

func (o *Node) String() string {
	b, _ := json.MarshalIndent(o, "", "    ")
	return string(b)
}

func (o *Node) PPrint() {
	fmt.Println(o.String())
}
