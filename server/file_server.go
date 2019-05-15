package server

import (
	"encoding/json"
	"net/http"
)

// http file server wrapper implementation allowing to create custom error messages (for status  >= 400)

// StatusRespWr ...
type StatusRespWr struct {
	http.ResponseWriter
	status int
}

// WriteHeader ...
func (w *StatusRespWr) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func customNotFoundWrapper(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srw := &StatusRespWr{ResponseWriter: w}
		h.ServeHTTP(srw, r)
		if srw.status >= 400 {
			json.NewEncoder(w).Encode("the requested image does not exist or an error occurred while processing, try loading image again")
		}
	}
}
