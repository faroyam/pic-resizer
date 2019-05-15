package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/faroyam/pic-resizer/resizer"

	b64 "encoding/base64"

	c "github.com/faroyam/pic-resizer/config"
	"github.com/faroyam/pic-resizer/msg"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

var cfg *c.Config

func init() {
	cfg = c.GetConfig()
	err := os.Mkdir(cfg.ImagesPath, os.ModePerm)
	if err != nil && os.IsNotExist(err) {
		log.Fatalf("can't create directory to save images: %v", err)
	}
}

// Serve starts http server
func Serve() error {

	r := mux.NewRouter().StrictSlash(true)

	s := http.StripPrefix(cfg.ImagesURL, http.FileServer(http.Dir(cfg.ImagesPath)))
	r.PathPrefix(cfg.ImagesURL).Handler(s)

	r.Handle(cfg.RootURL+cfg.GetURL,
		middlewareLog(http.HandlerFunc(getImageH))).Methods(http.MethodGet)

	r.Handle(cfg.RootURL+cfg.MultipartURL,
		middlewareLog(http.HandlerFunc(multipartH))).Methods(http.MethodPost)

	r.Handle(cfg.RootURL+cfg.Base64URL,
		middlewareLog(http.HandlerFunc(base64H))).Methods(http.MethodPost)

	r.Handle(cfg.RootURL,
		middlewareLog(http.HandlerFunc(rootH))).Methods(http.MethodGet)

	return http.ListenAndServe(cfg.Addr, r)
}

func rootH(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := msg.NewResponse(w, "it works"); err != nil {
		log.Println(err)
	}
}

func getImageH(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	url := r.URL.Query().Get("url")
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		msg.NewResponse(w, "fetching image error")
		return
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	img, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		msg.NewResponse(w, "preparing image error")
		return
	}

	s, err := sendImageToQueue(img)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		msg.NewResponse(w, "preparing image error")
		return
	}

	msg.NewResponse(w, fmt.Sprintf("your pic: %v%v.jpg, preview: %v%v%v.jpg", cfg.ImagesURL, s, cfg.ImagesURL, s, cfg.ImagesPreviewSuff))
}

func multipartH(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/form-data")
	msg.NewResponse(w, "pic from multipart/form-data")
}

func base64H(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	req := &msg.BaseRequest{}
	err := decoder.Decode(req)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		msg.NewResponse(w, "request format error")
		return
	}

	b, err := b64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		msg.NewResponse(w, "request data error")
		return
	}

	s, err := sendImageToQueue(b)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		msg.NewResponse(w, "preparing image error")
		return
	}

	msg.NewResponse(w, fmt.Sprintf("your pic: %v%v.jpg, preview: %v%v%v.jpg", cfg.ImagesURL, s, cfg.ImagesURL, s, cfg.ImagesPreviewSuff))
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v request to %v from %v", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func sendImageToQueue(b []byte) (string, error) {

	k := ksuid.New()
	s := k.String()
	resizer.Send(msg.Pic{ID: k.String(), Data: b})

	return s, nil
}
