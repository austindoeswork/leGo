package domain

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"time"

	"git.ottoq.com/otto-backend/valet/gen/namecase"
)

var BasePath = path.Join(os.Getenv("GOPATH"), "/src/git.ottoq.com/otto-backend/valet/domain")

var TypeIDOf = map[string]string{
	"Node": "0C74DFC158C646C280BCB0DAF9E015D1",
	"Desk": "E1874C161CDB492FB95EF210E653B886",
}

var List = []Object{ //ORDER MATTERS HERE VVV
	Object{
		Name:        namecase.New("Node"),
		Description: "Node represents a node in the organization permission heirarchy tree",
		TypeID:      TypeIDOf["Node"],
		Imports: []string{
			"time",
		},
		Parameters: []Parameter{
			ID(),
			TypeID(TypeIDOf["Node"]),
			Timestamp(),
			String("Name"),
		},
	},
	Object{
		Name:        namecase.New("Desk"),
		Description: "Desk where car keys can be stored",
		TypeID:      TypeIDOf["Desk"],
		Imports: []string{
			"time",
		},
		Parameters: []Parameter{
			ID(),
			TypeID(TypeIDOf["Desk"]),
			Timestamp(),
			String("Name"),
			Float("Lat"),
			Float("Lng"),
			ForeignK("NodeID", "Node", "ID"),
		},
	},
}

///////////////////
// HELPERS
///////////////////

func ID() Parameter {
	id := PrimaryK("ID")
	id.ConstructorOverride = "entity.UUID()"
	return id
}

func TypeID(typeid string) Parameter {
	tid := String("TypeID")
	tid.SQLType = "BINARY(16)"
	tid.ConstructorOverride = fmt.Sprintf(`"%s"`, typeid)
	return tid
}

func Timestamp() Parameter {
	t := Datetime("Timestamp")
	t.ConstructorOverride = "entity.Now()"
	return t
}

func PrimaryK(name string) Parameter {
	return Parameter{
		Name:       namecase.New(name),
		Type:       reflect.TypeOf(""),
		SQLType:    "BINARY(16)",
		Index:      true,
		PrimaryKey: true,
		ForeignKey: nil,
	}
}

func ForeignK(name, table, column string) Parameter {
	return Parameter{
		Name:       namecase.New(name),
		Type:       reflect.TypeOf(""),
		SQLType:    "BINARY(16)",
		Index:      true,
		PrimaryKey: false,
		ForeignKey: &ForeignKey{
			Table:  table,
			Column: column,
		},
	}
}

func String(name string) Parameter {
	return Parameter{
		Name:       namecase.New(name),
		Type:       reflect.TypeOf(""),
		SQLType:    "VARCHAR(100)",
		Index:      false,
		PrimaryKey: false,
		ForeignKey: nil,
	}
}

func Float(name string) Parameter {
	return Parameter{
		Name:       namecase.New(name),
		Type:       reflect.TypeOf(float64(0)),
		SQLType:    "FLOAT",
		Index:      false,
		PrimaryKey: false,
		ForeignKey: nil,
	}
}

func Datetime(name string) Parameter {
	return Parameter{
		Name:       namecase.New(name),
		Type:       reflect.TypeOf(time.Now()),
		SQLType:    "DATETIME",
		Index:      false,
		PrimaryKey: false,
		ForeignKey: nil,
	}
}
