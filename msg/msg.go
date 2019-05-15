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

// BaseRequest decribes message with base64 string
type BaseRequest struct {
	Data string `json:"data"`
}

type response struct {
	Service string `json:"service"`
	Comment string `json:"answer"`
	Images  struct {
		Original string `json:"original_image"`
		Resized  string `json:"resized_image"`
	} `json:"images"`
}

// NewResponse writes json to http.ResponseWriter
func NewResponse(w http.ResponseWriter, s, originalURL, resizedURL string) error {
	resp := response{Service: cfg.ServiceName, Comment: s}
	resp.Images.Original = originalURL
	resp.Images.Resized = resizedURL
	return json.NewEncoder(w).Encode(resp)
}
