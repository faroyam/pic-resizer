package msg

import (
	"encoding/json"
	"net/http"

	c "github.com/faroyam/pic-resizer/config"
)

var cfg *c.Config

func init() {
	cfg = c.GetConfig()
}

// Pic decribes message from http server to image resize service
type Pic struct {
	ID   string
	Data []byte
}

type response struct {
	Service string `json:"service"`
	Answer  string `json:"answer"`
}

// NewResponse write json to http.ResponseWriter
func NewResponse(w http.ResponseWriter, s string) error {
	var resp = response{Service: cfg.ServiceName, Answer: s}
	return json.NewEncoder(w).Encode(resp)
}

// BaseRequest decribes message with base64 string
type BaseRequest struct {
	Data string `json:"data"`
}
