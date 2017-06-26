package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"git.ottoq.com/otto-backend/valet/gen/database"
	"git.ottoq.com/otto-backend/valet/gen/domain"
)

var (
	rundir  = path.Join(os.Getenv("GOPATH"), "/src/git.ottoq.com/otto-backend/valet")
	funcMap = template.FuncMap{ //added func to compare strings
		"contains": func(a, b string) bool {
			return strings.Contains(a, b)
		},
	}
)

func main() {
	verifyRunDir()
	Domain()
	Database()
}

func Domain() error {
	basepath := domain.BasePath
	for _, d := range domain.List {
		MakePackage(path.Join(basepath, d.Name.Lower), d.Name.Lower+"_gen.go", "Domain", domain.Plate["Domain"], d)
	}
	return nil
}

func Database() error {
	//Database is also generated from domain objects
	basepath := database.BasePath
	MakePackage(basepath, "database_gen.go", "Database", database.Plate["Database"], domain.List)
	return nil
}

func MakePackage(basedir, filename, tmplname, tmpl string, obj interface{}) error {
	code := GenerateCode(tmplname, tmpl, obj)
	if err := os.MkdirAll(basedir, 0744); err != nil {
		return err
	}
	f, err := os.Create(path.Join(basedir, filename))
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(code)); err != nil {
		return err
	}
	return nil
}

func GenerateCode(tmplname, tmpl string, obj interface{}) string {
	t, err := template.New(tmplname).Funcs(funcMap).Parse(tmpl)
	if err != nil {
		log.Fatal(err.Error())
	}
	b := new(bytes.Buffer)
	if err := t.Execute(b, obj); err != nil {
		log.Fatal(err.Error())
	}
	formatted, err := format.Source(b.Bytes())
	if err != nil {
		fmt.Printf("%s\n", b.Bytes())
		log.Fatal(err.Error())
	}

	return string(formatted)
}

func verifyRunDir() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get current working directory")
	}
	if cwd != rundir {
		log.Fatalf("must be run from directory: %s\n", rundir)
	}
}
