package inputsample

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/wardn/uuid"

	"git.ottoq.com/otto-backend/valet/server"
	"git.ottoq.com/otto-backend/valet/server/securecookie"
	"git.ottoq.com/otto-backend/valet/server/session"
)

const (
	TypeID = "2BE17DF8BBCD43FB8FE53811AED9986D"
)

type Payload struct {
	id       string
	w        http.ResponseWriter
	r        *http.Request
	sesh     string
	Contents Contents
}

func (p *Payload) Writer() http.ResponseWriter {
	return p.w
}
func (p *Payload) Request() *http.Request {
	return p.r
}
func (p *Payload) SessionID() string {
	return p.sesh
}
func (p *Payload) ID() string {
	return p.id
}
func (p *Payload) TypeID() string {
	return TypeID
}

type Contents struct {
	Name   string
	Number int
}

// FromHTTPRequest takes an http request/response and returns a Request.
func FromHTTPRequest(w http.ResponseWriter, r *http.Request,
	sc *securecookie.Config,
	sessionCookieName string) (server.InputDTO, error) {

	var c Contents
	seshID, err := ParseRW(&c, w, r, sc, sessionCookieName)
	if err != nil {
		return nil, err
	}
	return &Payload{
		id:       uuid.NewNoDash(),
		w:        w,
		r:        r,
		sesh:     seshID,
		Contents: c,
	}, nil
}

// Parses object and returns a sessionID
func ParseRW(
	obj interface{},
	w http.ResponseWriter, r *http.Request,
	sc *securecookie.Config, sessionCookieName string) (string, error) {

	if r == nil {
		return "", fmt.Errorf("NIL REQUEST")
	}
	cookies := r.Cookies()

	// decode the session info from the cookies
	now := time.Now().UTC()
	sesh, err := session.FromCookies(cookies, sessionCookieName, sc)
	if err != nil {
		return "", err
	}
	if sesh == nil {
		return "", fmt.Errorf("NIL SESSION")
	}
	if sesh.Timestamp.After(now) {
		if err := session.SetCookie(sesh, sessionCookieName, w, sc); err != nil {
			return "", err
		}
	}

	// unmarshall the request into the object interface
	raw := requestBody(1024, w, r)

	if len(raw) == 0 {
		return sesh.ID, nil
	}
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		return "", err
	}
	return sesh.ID, nil
}

func requestBody(maxBodySize int64, w http.ResponseWriter, r *http.Request) string {
	if r == nil || r.Body == nil {
		return ""
	}
	// enforce a size limit to prevent abuse
	mbr := http.MaxBytesReader(w, r.Body, maxBodySize)

	// read the bytes
	b, err := ioutil.ReadAll(mbr)
	if err != nil {
		return ""
	}
	return string(b)
}

//
