package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"

	"git.ottoq.com/otto-backend/valet/config"
	"git.ottoq.com/otto-backend/valet/database"
	"git.ottoq.com/otto-backend/valet/dto/input/sample"
	"git.ottoq.com/otto-backend/valet/entity"
	"git.ottoq.com/otto-backend/valet/server"
	"git.ottoq.com/otto-backend/valet/server/securecookie"
)

const (
	appName           = "valet"
	cookieDomain      = "ottoq.com"
	cookiePath        = "/"
	cookieSessionName = "v"
	configFileName    = "config"
	serverIDFileName  = "serverid"
)

var (
	// places to look (in order) for the config file
	configPaths = []string{
		"./",
		path.Join(homeDir(), ".config/ottoq/"+appName),
		"/usr/local/etc/ottoq/" + appName,
		"/etc/ottoq/" + appName,
	}
	cookieMinAge = int64(0)
	cookieMaxAge = int64(10 * 365 * 24 * 60 * 60) // 10 years... because yeah
)

////////////////////////////////////////////////////////////

// homeDir gets the user's home directory
func homeDir() string {
	u, err := user.Current()
	if err != nil {
		return ""
	}
	return u.HomeDir
}

////////////////////////////////////////////////////////////

// getConfig returns a config file if found, otherwise an error
func getConfig() (*config.Config, error) {
	for _, p := range configPaths {
		configPath := path.Join(p, configFileName)
		configPath = path.Clean(configPath)
		b, err := ioutil.ReadFile(configPath)
		if err != nil {
			continue
		}
		var c config.Config
		if err := json.Unmarshal(b, &c); err != nil {
			continue
		}
		if err := c.Validate(); err != nil {
			continue
		}
		return &c, nil
	}
	return nil, errors.New("failed to read config")
}

////////////////////////////////////////////////////////////

func main() {
	var err error

	// CONFIG
	c, err := getConfig()
	if err != nil {
		log.Fatalf("Fatal: Config error: %s\nExample: %s\n", err.Error(), config.Example())
	}

	// DATABASE
	dbadd := "tcp(127.0.0.1:3306)"
	db, err := database.New(dbadd, "test3", "austin", "")
	if err != nil {
		log.Fatal(err)
	}
	_ = db

	// SECURE COOKIE
	cc, err := securecookie.New(c.HashKey(), c.BlockKey(),
		cookieDomain, cookiePath,
		cookieMinAge, cookieMaxAge,
		c.Secure(), c.Secure())
	if err != nil {
		log.Fatalf("Fatal: Failed to initialize cookie manager. Error: %s\n", err.Error())
	}

	// SERVER
	s, err := server.New(c.Secure(), c.InsecureRedirect(), c.Gzip(), c.Cache(),
		c.ServerAddress(), cookieDomain, cookieSessionName, certPath(), keyPath(),
		cc)

	s.RegisterHTTPRoute("/test", server.HTTPConverterMap{"POST": inputsample.FromHTTPRequest})
	s.RegisterHandler(H{})

	s.Start()
}

type H struct {
}

func (h H) InputTypeID() string {
	return inputsample.TypeID
}

func (h H) Notify(input server.InputDTO, response chan entity.Identifier) error {
	// input.Writer().Write([]byte("AYYYLMAO"))
	input.Writer().Write([]byte("hello world."))
	in := input.(*inputsample.Payload)
	fmt.Println(in.Contents.Name)
	response <- nil
	return nil
}

func pemPath(name string) string {
	for _, p := range configPaths {
		p := path.Join(p, name)
		p = path.Clean(p)
		// ensure there's currently a file here
		_, err := os.Open(p)
		if err != nil {
			continue
		}
		return p
	}
	return ""
}

func certPath() string {
	return pemPath("cert.pem")
}

func keyPath() string {
	return pemPath("key.pem")
}

//
