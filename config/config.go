/*

Package config is used to define and store various server configuration
parameters.

*/
package config

import (
	"encoding/base32"
	"encoding/json"

	"git.ottoq.com/otto-backend/valet/server/securecookie"
)

////////////////////////////////////////////////////////////

// ErrInvalidConfig is an invalid config error
type ErrInvalidConfig struct {
	Reason string
}

// Error returns the error string
func (err *ErrInvalidConfig) Error() string {
	return err.Reason
}

////////////////////////////////////////////////////////////

// Config stores server configuration
type Config struct {
	*config
}

type config struct {
	ServerAddress    string // ServerAddress is ths address that this server starts up as.
	DatabaseAddress  string // DatabaseAddress is the address of the db server.
	DatabaseName     string // DatabaseName is the name of the database to use.
	LogFilePath      string // LogFilePath is the path for the log file.
	HashKey          string // HashKey is for verifying cookie integrity.
	BlockKey         string // BlockKey is for encrypting cookies.
	Secure           bool   // Secure determines http/https, ws/wss, secure cookies etc.
	InsecureRedirect bool   // InsecureRedirect will redirect port 80 calls if enabled.
	Gzip             bool   // Gzip indicates if gzip compression should be used.
	Cache            bool   // Cache indicates if data should be stored in memory after compression.
}

////////////////////////////////////////////////////////////

func generateRandomKey() string {
	return base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(16))
}

// Example returns an example config
func Example() string {
	c := &Config{
		&config{
			ServerAddress:    "localhost:8080",
			DatabaseAddress:  "devdb.ottoq.com:28015",
			DatabaseName:     "ottoq",
			LogFilePath:      "/var/log/ottoq/ottoqvalet/log",
			HashKey:          generateRandomKey(),
			BlockKey:         generateRandomKey(),
			Secure:           true,
			InsecureRedirect: true,
			Gzip:             true,
			Cache:            true,
		},
	}
	b, _ := json.MarshalIndent(&c, "", "  ")
	return string(b)
}

// Validate validates a config
func (c *Config) Validate() error {
	if c.config == nil {
		return &ErrInvalidConfig{"config is empty"}
	}
	if len(c.config.ServerAddress) == 0 {
		return &ErrInvalidConfig{"unspecified server address"}
	}
	if len(c.config.DatabaseAddress) == 0 {
		return &ErrInvalidConfig{"unspecified database address"}
	}
	if len(c.config.DatabaseName) == 0 {
		return &ErrInvalidConfig{"unspecified database name"}
	}
	if len(c.config.LogFilePath) == 0 {
		return &ErrInvalidConfig{"unspecified log file path"}
	}
	if len(c.config.HashKey) == 0 {
		return &ErrInvalidConfig{"unspecified hash key"}
	}
	if len(c.config.BlockKey) == 0 {
		return &ErrInvalidConfig{"unspecified block key"}
	}
	return nil
}

////////////////////////////////////////////////////////////

// ServerAddress returns the server address
func (c *Config) ServerAddress() string {
	if c.config == nil {
		return ""
	}
	return c.config.ServerAddress
}

// DatabaseAddress returns the database address
func (c *Config) DatabaseAddress() string {
	if c.config == nil {
		return ""
	}
	return c.config.DatabaseAddress
}

// DatabaseName returns the database name
func (c *Config) DatabaseName() string {
	if c.config == nil {
		return ""
	}
	return c.config.DatabaseName
}

// LogFilePath returns the log file path
func (c *Config) LogFilePath() string {
	if c.config == nil {
		return ""
	}
	return c.config.LogFilePath
}

// HashKey returns the hash key used for encoding/decoding cookies
func (c *Config) HashKey() []byte {
	if c.config == nil {
		return nil
	}
	b, err := base32.StdEncoding.DecodeString(c.config.HashKey)
	if err != nil {
		return nil
	}
	return b
}

// BlockKey returns the block key used for encrypting/decrypting cookies
func (c *Config) BlockKey() []byte {
	if c.config == nil {
		return nil
	}
	b, err := base32.StdEncoding.DecodeString(c.config.BlockKey)
	if err != nil {
		return nil
	}
	return b
}

// Secure returns whether we're running in secure mode
//
// This will determine the protocols we use and the cookie security settings
func (c *Config) Secure() bool {
	if c.config == nil {
		// default restricted
		return true
	}
	return c.config.Secure
}

// InsecureRedirect returns whether the insecure redirect is enabled
func (c *Config) InsecureRedirect() bool {
	if c.config == nil {
		// default restricted
		return false
	}
	return c.config.InsecureRedirect
}

// Gzip returns whether to gzip compress response data.
func (c *Config) Gzip() bool {
	if c.config == nil {
		return false
	}
	return c.config.Gzip
}

// Cache returns whether to cache compressed data.
func (c *Config) Cache() bool {
	return false
}

////////////////////////////////////////////////////////////
