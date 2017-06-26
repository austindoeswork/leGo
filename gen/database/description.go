package database

import (
	"os"
	"path"
)

var BasePath = path.Join(os.Getenv("GOPATH"), "/src/git.ottoq.com/otto-backend/valet/database")
