package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"git.ottoq.com/otto-backend/valet/entity"
	"git.ottoq.com/otto-backend/valet/server/securecookie"
)

// Handler represents instances that are interested in being notified
// when request of a given type are received
type Handler interface {
	InputTypeID() string
	// OutputTypeID() map[string]struct{}
	// Notify is the method used to comminucate with the handler
	Notify(input InputDTO, responses chan entity.Identifier) error
}

// HandlerMap associates request types with handlers
type handlerMap map[string]Handler

type Server struct {
	*http.Server

	secure            bool
	insecureRedirect  bool
	gzip              bool
	cache             bool
	tlsCert           string
	tlsKey            string
	cookieDomain      string
	sessionCookieName string

	// logger            entity.Logger
	secureCookie *securecookie.Config
	mux          *http.ServeMux
	handlers     handlerMap
	// sockets  *socketMap
}

func New(
	secure, insecureRedirect, gzip, cache bool,
	addr, cookieDomain, sessionCookieName, tlsCert, tlsKey string,
	// logger entity.Logger,
	sc *securecookie.Config,
) (*Server, error) {

	if sc == nil {
		return nil, fmt.Errorf("NIL COOKIE CONFIG")
	}
	if secure && (len(tlsCert) == 0 || len(tlsKey) == 0) {
		return nil, fmt.Errorf("SECURE CERT NOT FOUND")
	}

	// create our server
	s := &Server{
		Server: &http.Server{
			Addr: addr,
		},
		cookieDomain:      cookieDomain,
		sessionCookieName: sessionCookieName,
		secure:            secure,
		insecureRedirect:  insecureRedirect,
		gzip:              gzip,
		cache:             cache,
		tlsCert:           tlsCert,
		tlsKey:            tlsKey,
		mux:               http.NewServeMux(),
		// logger:            logger,
		secureCookie: sc,
		handlers:     handlerMap{},
		// sockets:      &socketMap{s: make(map[string]struct{})},
	}

	// s.mux.Handle(UserLogoutPath, s.logout())
	return s, nil
}

type InputDTO interface {
	TypeID() string
	Writer() http.ResponseWriter
	Request() *http.Request
	SessionID() string
}

// Start initiates the server's listener
func (s *Server) Start() {
	log.Printf(
		"Starting server...\nSecure: %t\nInsecureRedirect: %t\nAddress: %s\n",
		s.secure,
		s.insecureRedirect,
		s.Addr,
	)
	// Use TLS if secure is true
	if s.secure {
		var wg sync.WaitGroup
		wg.Add(1)
		if s.insecureRedirect {
			go func() {
				mux := http.NewServeMux()
				mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
					host := req.Host
					if strings.Contains(host, ":") {
						h, _, err := net.SplitHostPort(req.Host)
						if err != nil {
							// s.logger.Log(genericerror.New("03e41d2b-b000-44e2-8e2d-b7e761ae1279", err))
						}
						host = h
					}
					_, primaryPort, err := net.SplitHostPort(s.Addr)
					if err != nil {
						// s.logger.Log(genericerror.New("453ba0d5-018f-42b2-8630-0aced5bbfffb", err))
					}
					redirect, err := url.Parse(fmt.Sprintf("https://%s:%s%s", host, primaryPort, req.URL.String()))
					if err != nil {
						// s.logger.Log(genericerror.New("d56369d2-5ef9-4904-8350-b8eeb56254f7", err))
					}
					http.Redirect(w, req, redirect.String(), http.StatusSeeOther)
				})
				if err := http.ListenAndServe(":80", mux); err != nil {
					log.Println(err.Error())
					// s.logger.Log(genericerror.New("9c76d39f-6179-4504-95ed-1e2a226257a2", err))
				}
				wg.Done()
			}()
		}
		go func() {
			if err := http.ListenAndServeTLS(s.Addr, s.tlsCert, s.tlsKey, s.mux); err != nil {
				log.Println(err.Error())
				// s.logger.Log(genericerror.New("bd3a9017-b003-4714-be7f-05836f18b784", err))
			}
			wg.Done()
		}()
		wg.Wait()
		os.Exit(0)
	}
	// Not secure, so don't use TLS
	if err := http.ListenAndServe(s.Addr, s.mux); err != nil {
		// s.logger.Log(genericerror.New("1b2947f0-7ce4-4600-80cc-b63ede0789a2", err))
		os.Exit(0)
	}
}

// HTTPRequestToInput is a function that converts an http request into an http input
type HTTPRequestToInput func(w http.ResponseWriter, r *http.Request,
	sc *securecookie.Config, sessionCookieName string) (InputDTO, error)

// HTTPConverterMap maps http routes (GET/POST) to a function that converts it
type HTTPConverterMap map[string]HTTPRequestToInput

// RegisterHTTPRoute registers http routes along with their route methods for handling entities
func (s *Server) RegisterHTTPRoute(route string, methods HTTPConverterMap) {
	s.mux.Handle(route, s.httpRequestHandler(methods))
}

// RegisterHandler accepts a handler and registers it in the HandlerMap
func (s *Server) RegisterHandler(h Handler) {
	s.handlers[h.InputTypeID()] = h
}
