package session

import (
	"encoding/base64"
	"net/http"
	"time"

	"git.ottoq.com/otto-backend/valet/server/securecookie"

	"github.com/wardn/uuid"
)

// Session is the interface for sessions
type Session struct {
	ID        string
	Token     string
	Timestamp time.Time
}

// FromCookies returns a session if it is already set, or a new session if one doesn't exist.
func FromCookies(cookies []*http.Cookie, cookieName string,
	sc *securecookie.Config) (*Session, error) {

	sesh := sessionFromCookies(cookies, cookieName, sc)
	if sesh != nil {
		return sesh, nil
	}

	// This token is used for CSRF mitigation
	token := base64.RawStdEncoding.EncodeToString([]byte(uuid.NewNoDash()))
	return &Session{
		ID:        uuid.NewNoDash(),
		Token:     token,
		Timestamp: time.Now(),
	}, nil
}

func sessionFromCookies(cookies []*http.Cookie, cookieName string, sc *securecookie.Config) *Session {
	for _, c := range cookies {
		if c.Name != cookieName {
			continue
		}
		var s Session
		if err := sc.Decode(c.Name, c.Value, &s); err != nil {
			continue
		}
		return &s
	}
	return nil
}

// SetCookie sets a session cookie.
func SetCookie(session *Session, cookieName string, w http.ResponseWriter, sc *securecookie.Config) error {
	if w == nil {
		return nil
	}
	value, err := sc.Encode(cookieName, session)
	if err != nil {
		return err
	}
	c := &http.Cookie{
		Name:     cookieName,
		Value:    value,
		Domain:   sc.Domain(),
		Path:     sc.Path(),
		MaxAge:   sc.MaxAge(),
		Secure:   sc.Secure(),
		HttpOnly: sc.HTTPOnly(),
	}
	http.SetCookie(w, c)
	return nil
}

//
