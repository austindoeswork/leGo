package domain

import (
	"reflect"
	"strings"

	"git.ottoq.com/otto-backend/valet/gen/namecase"
)

type Object struct {
	Name        *namecase.Name
	Description string
	TypeID      string
	Imports     []string
	Parameters  []Parameter
}

type Parameter struct {
	Name                *namecase.Name
	Type                reflect.Type
	SQLType             string
	Index               bool
	PrimaryKey          bool
	ForeignKey          *ForeignKey
	ConstructorOverride string
}

type ForeignKey struct {
	Table  string
	Column string
}

func (o Object) SQLInsert() string {
	params := []string{}
	for _, p := range o.Parameters {
		if strings.Contains(p.Name.UpperCamel, "ID") {
			params = append(params, `UNHEX( '%s' )`)
		} else {
			switch p.Type.String() {
			case "float64":
				params = append(params, `%f`)
			case "int":
				params = append(params, `%d`)
			case "string":
				params = append(params, `'%s'`)
			case "time.Time":
				params = append(params, `'%s'`)
			}
		}
	}
	paramstr := strings.Join(params, ",\n")
	insertstr := "INSERT INTO " + o.Name.UpperCamel + " VALUES(\n" +
		paramstr + "\n);"
	return insertstr
}

func (o Object) SQLSchema() string {
	columns := []string{}
	primary := []string{}
	secondary := []string{}
	for _, p := range o.Parameters {
		columns = append(columns, p.Name.UpperCamel+" "+p.SQLType)
		if p.PrimaryKey {
			primary = append(primary,
				PrimaryString(p.Name.UpperCamel))
		}
		if p.ForeignKey != nil {
			secondary = append(secondary,
				ForeignString(p.Name.UpperCamel, p.ForeignKey.Table, p.ForeignKey.Column))
		}
	}
	columns = append(columns, primary...)
	columns = append(columns, secondary...)

	colstr := strings.Join(columns, ",\n")
	tablstr := "CREATE TABLE " + o.Name.UpperCamel + " (\n" +
		colstr + "\n);"

	return tablstr
}

func PrimaryString(columns ...string) string {
	colstr := strings.Join(columns, ", ")
	primstr := "PRIMARY KEY (" + colstr + ")"
	return primstr
}
func ForeignString(localCol, table, foreigncol string) string {
	forstr := "FOREIGN KEY (" + localCol + ")" + " REFERENCES " + table + "(" + foreigncol + ")"
	return forstr
}
