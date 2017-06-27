// Package securecookie uses AES encryption to encrypt cookie content along
// with a keyed-hash message authentication code (HMAC) to ensure data
// integrity
package securecookie

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"hash"
	"io"
	"math"
	"strconv"
	"time"
)

// Config provides convenience methods for storing and retrieving secure cookie values
type Config struct {
	hashKey    []byte           // hashKey is the key used for generating a hash of the cookie
	hashFunc   func() hash.Hash // hashFunc is the function used for hashing the cookie
	blockKey   []byte           // blockKey is used for encrypting the cookie value
	block      cipher.Block     // block is the cipher used to encrypt the cookie value
	domain     string           // domain defines the cookie's domain scope
	path       string           // path defines the cookie's path scope
	minAge     int64            // minAge is the minimum acceptable age of the cookie in seconds
	maxAge     int64            // maxAge is the maximum acceptable age of the cookie in seconds
	maxLength  int              // maxLength is the maximum length of the encoded cookie
	secure     bool             // secure cookies are limited to secure channels (defined by the user agent)
	httpOnly   bool             // httpOnly cookies have a scope limited to HTTP requests
	serializer Serializer       // serializer is the serializer used to encode/decode the cookie
}

// Serializer provides methods to serialize and deserialize values
type Serializer interface {
	Serialize(src interface{}) ([]byte, error)
	Deserialize(src []byte, dst interface{}) error
}

// New instantiates and returns a new cookie instance or returns an error
func New(hashKey, blockKey []byte, domain, path string, minAge, maxAge int64, secure, httpOnly bool) (*Config, error) {
	if hashKey == nil {
		return nil, new(ErrHashKeyNotSet)
	}
	if blockKey == nil {
		return nil, new(ErrBlockKeyNotSet)
	}
	block, err := aes.NewCipher(blockKey)
	if err != nil {
		return nil, err
	}
	return &Config{
		hashKey:    hashKey,
		hashFunc:   sha256.New,
		blockKey:   blockKey,
		block:      block,
		domain:     domain,
		path:       path,
		minAge:     minAge,
		maxAge:     maxAge,
		secure:     secure,
		httpOnly:   httpOnly,
		maxLength:  4096,
		serializer: new(GobEncoder),
	}, nil
}

// Domain returns the cookie domain.
func (c *Config) Domain() string {
	return c.domain
}

// Path returns the cookie path.
func (c *Config) Path() string {
	return c.path
}

// MinAge returns the cookie's minimum age
func (c *Config) MinAge() int {
	return ageToInt(c.minAge)
}

// MaxAge returns the cookie's maximum age
func (c *Config) MaxAge() int {
	return ageToInt(c.maxAge)
}

func ageToInt(i int64) int {
	if i < 0 {
		return -1
	}
	if i > math.MaxInt32 {
		return math.MaxInt32
	}
	return int(i)
}

// Secure returns whether the cookie will be secure
func (c *Config) Secure() bool {
	return c.secure
}

// HTTPOnly returns whether the cookie will be httpOnly
func (c *Config) HTTPOnly() bool {
	return c.httpOnly
}

// Encode serializes, encrypts, signs, and encodes a cookie value.
//
// The name argument is the cookie name.  It is stored with the encoded value.
// The argument is the value to be encoded.  It can be any value that can be
// encoded with the given serializer.
func (c *Config) Encode(name string, value interface{}) (string, error) {
	b, err := c.serializer.Serialize(value)
	if err != nil {
		return "", err
	}
	b, err = encrypt(c.block, b)
	if err != nil {
		return "", err
	}
	b = encode(b)
	b = []byte(fmt.Sprintf("%s|%d|%s|", name, time.Now().UTC().Unix(), b))
	mac := createMac(hmac.New(c.hashFunc, c.hashKey), b[:len(b)-1])
	b = append(b, mac...)[len(name)+1:]
	b = encode(b)
	if c.maxLength != 0 && len(b) > c.maxLength {
		return "", new(ErrEncodedLengthTooLong)
	}
	return string(b), nil
}

// Decode decodes, verifies the mac, decrypts and deserializes a cookie value.
//
// The name argument is the cookie name and must be the same name that was used
// when it was stored.  The value argument is the encoded cookie value. The dst
// argument is where the cookie will be decoded and must be a pointer.
func (c *Config) Decode(name, value string, dst interface{}) error {
	if c.maxLength != 0 && len(value) > c.maxLength {
		return new(ErrDecodeValueTooLong)
	}
	b, err := decode([]byte(value))
	if err != nil {
		return err
	}
	parts := bytes.SplitN(b, []byte("|"), 3)
	if len(parts) != 3 {
		return new(ErrInvalidMac)
	}
	h := hmac.New(c.hashFunc, c.hashKey)
	b = append([]byte(name+"|"), b[:len(b)-len(parts[2])-1]...)
	if err = verifyMac(h, b, parts[2]); err != nil {
		return err
	}
	t1, err := strconv.ParseInt(string(parts[0]), 10, 64)
	if err != nil {
		return err
	}
	t2 := time.Now().UTC().Unix()
	if c.minAge != 0 && t1 > t2-c.minAge {
		return new(ErrTimestampTooNew)
	}
	if c.maxAge != 0 && t1 < t2-c.maxAge {
		return new(ErrTimestampExpired)
	}
	b, err = decode(parts[1])
	if err != nil {
		return err
	}
	if c.block != nil {
		if b, err = decrypt(c.block, b); err != nil {
			return err
		}
	}
	if err = c.serializer.Deserialize(b, dst); err != nil {
		return err
	}
	return nil
}

// GenerateRandomKey returns a random key of the given length
func GenerateRandomKey(length int) []byte {
	k := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return nil
	}
	return k
}

func encrypt(block cipher.Block, value []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	ciphertext := make([]byte, blockSize+len(value))
	iv := ciphertext[:blockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[blockSize:], value)
	return ciphertext, nil
}

func decrypt(block cipher.Block, value []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	plaintextSize := len(value) - blockSize
	if plaintextSize < 0 {
		return nil, new(ErrInvalidEncryptedValueLength)
	}
	plaintext := make([]byte, plaintextSize)
	stream := cipher.NewCTR(block, value[:blockSize])
	stream.XORKeyStream(plaintext, value[blockSize:])
	return plaintext, nil
}

func encode(value []byte) []byte {
	return []byte(base64.URLEncoding.EncodeToString(value))
}

func decode(value []byte) ([]byte, error) {
	return base64.URLEncoding.DecodeString(string(value))
}

func createMac(h hash.Hash, value []byte) []byte {
	h.Write(value)
	return h.Sum(nil)
}

func verifyMac(h hash.Hash, value []byte, mac []byte) error {
	mac2 := createMac(h, value)
	if len(mac) == len(mac2) && subtle.ConstantTimeCompare(mac, mac2) == 1 {
		return nil
	}
	return new(ErrInvalidMac)
}

// GobEncoder is a gob encoder.
type GobEncoder struct{}

// Serialize takes a value src and returns either the gob-encoded bytes,
// or an error
func (g *GobEncoder) Serialize(src interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(src); err != nil {
		return nil, new(ErrGobEncode)
	}
	return buf.Bytes(), nil
}

// Deserialize takes gob-encoded bytes and a destination type and deserializes
// the src into the dst. If an error is encountered it will be returned.
func (g *GobEncoder) Deserialize(src []byte, dst interface{}) error {
	dec := gob.NewDecoder(bytes.NewBuffer(src))
	if err := dec.Decode(dst); err != nil {
		return new(ErrGobDecode)
	}
	return nil
}

// ErrHashKeyNotSet is an error that results from a hash key not being set
type ErrHashKeyNotSet struct{}

// Error returns the error string
func (e *ErrHashKeyNotSet) Error() string { return "hash key not set" }

// ErrBlockKeyNotSet is an error that results from a block key not being set
type ErrBlockKeyNotSet struct{}

// Error returns the error string
func (e *ErrBlockKeyNotSet) Error() string { return "block key not set" }

// ErrGobEncode is an error that results from a failure to encode a gob
type ErrGobEncode struct{}

// Error returns the error string
func (e *ErrGobEncode) Error() string { return "failed to encode gob" }

// ErrGobDecode is an error that results from a failure to decode a gob
type ErrGobDecode struct{}

// Error returns the error string
func (e *ErrGobDecode) Error() string { return "failed to decode gob" }

// ErrInvalidEncryptedValueLength is an error that results when the length
// of the encrypted value is too short
type ErrInvalidEncryptedValueLength struct{}

// Error returns the error string
func (e *ErrInvalidEncryptedValueLength) Error() string {
	return "length of encrypted value is too short"
}

// ErrEncodedLengthTooLong is an error that results from an encoded length
// that's too long
type ErrEncodedLengthTooLong struct{}

// Error returns the error string
func (e *ErrEncodedLengthTooLong) Error() string {
	return "encoded length is too long"
}

// ErrDecodeValueTooLong is an error that results from a decode value
// that's too long
type ErrDecodeValueTooLong struct{}

// Error returns the error string
func (e *ErrDecodeValueTooLong) Error() string {
	return "decode value is too long"
}

// ErrInvalidMac is an error that results from an invalid mac
type ErrInvalidMac struct{}

// Error returns the error string
func (e *ErrInvalidMac) Error() string {
	return "mac is invalid"
}

// ErrTimestampTooNew is an error that results from a timestamp that's too new
type ErrTimestampTooNew struct{}

// Error returns the error string
func (e *ErrTimestampTooNew) Error() string {
	return "timestamp is too new"
}

// ErrTimestampExpired is an error that results from a timestamp that's expired
type ErrTimestampExpired struct{}

// Error returns the error string
func (e *ErrTimestampExpired) Error() string {
	return "timestamp is expired"
}

// ErrNilEncodeConfigValue is an error that results from a nil config value
// passed to encode.
type ErrNilEncodeConfigValue struct{}

// Error returns the error string
func (e *ErrNilEncodeConfigValue) Error() string {
	return "nil config value"
}
