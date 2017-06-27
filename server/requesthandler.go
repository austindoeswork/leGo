package server

import (
	"encoding/json"
	"net/http"
	"time"

	"git.ottoq.com/otto-backend/valet/entity"
)

type httpStatusHandler struct {
	// The "int" is the http status returned to the user.
	fn func(http.ResponseWriter, *http.Request) int
}

func (h *httpStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := h.fn(w, r)
	if status >= 400 {
		http.Error(w, http.StatusText(status), status)
	}
}

func (s *Server) httpRequestHandler(rm HTTPConverterMap) *httpStatusHandler {
	return &httpStatusHandler{
		fn: func(w http.ResponseWriter, r *http.Request) int {

			////////
			// Parse the input
			////////

			f, ok := rm[r.Method]
			if !ok {
				// s.logger.Log
				return http.StatusMethodNotAllowed
			}
			input, err := f(w, r, s.secureCookie, s.sessionCookieName)
			if err != nil {
				// s.logger.Log
				return http.StatusBadRequest
			}
			// s.logger.Log(req)

			////////
			// Let the handler know we have his request
			////////

			handler, ok := s.handlers[input.TypeID()]
			if !ok {
				// s.logger.Log
				return http.StatusInternalServerError
			}
			responses := make(chan entity.Identifier)

			go func(handler Handler, responses chan entity.Identifier) {
				err := handler.Notify(input, responses)
				if err != nil {
					// s.logger.Log
				}
			}(handler, responses)

			////////
			// Select on the handler's response
			////////

			select {
			case resp := <-responses:
				if resp == nil {
					break
				}
				// don't bypass the same-origin policy restriction (CORS, CSRF)
				// w.Header().Set("Access-Control-Allow-Origin", "*")
				// prevent clickjacking
				w.Header().Set("X-FRAME-OPTIONS", "DENY")
				if err := json.NewEncoder(w).Encode(resp); err != nil {
					// s.logger.Log
					return http.StatusInternalServerError
				}
			case <-time.After(5 * time.Second):
				return http.StatusInternalServerError
			}
			return http.StatusOK
		},
	}
}
