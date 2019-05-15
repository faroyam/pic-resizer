package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var doOnce sync.Once
var cfg *Config

// Config stores app startup options
type Config struct {
	ServiceName       string `json:"service_name"`
	Addr              string `json:"addr"`
	RootURL           string `json:"root_url"`
	GetURL            string `json:"get_url"`
	MultipartURL      string `json:"multipart_url"`
	Base64URL         string `json:"base_url"`
	ImagesURL         string `json:"image_url"`
	ImagesPath        string `json:"image_path"`
	ImagesPreviewSuff string `json:"image_preview_suffix"`
}

// GetConfig read config file on startup and returns configeration
func GetConfig() *Config {
	doOnce.Do(func() {
		cfg = &Config{
			ServiceName:       "Pic-resizer",
			Addr:              "0.0.0.0:8080",
			RootURL:           "/v1",
			GetURL:            "/get",
			MultipartURL:      "/multipart",
			Base64URL:         "/base",
			ImagesURL:         "/images/",
			ImagesPath:        "./images/",
			ImagesPreviewSuff: "_preview",
		}

		defer func() {
			log.Printf(`%v started...
	guide:
	to send image via link send GET request to: %v%v%v?url=http://example.com/image.jpg
	to send image as base64 string send POST request {"data":"YXNk..."} to: %v%v%v
	to send image as multipart/form-data send POST request to: %v%v%v
	resized image will appear at: %v%v<unique string>%v.jpg
	original image will be stored at: %v%v<unique string>.jpg
	images directory: %v`,
				cfg.ServiceName,
				cfg.Addr, cfg.RootURL, cfg.GetURL,
				cfg.Addr, cfg.RootURL, cfg.Base64URL,
				cfg.Addr, cfg.RootURL, cfg.MultipartURL,
				cfg.Addr, cfg.ImagesURL, cfg.ImagesPreviewSuff,
				cfg.Addr, cfg.ImagesURL,
				cfg.ImagesPath,
			)
		}()

		jsonFile, err := os.Open("config.json")
		if err != nil {
			log.Println(err)
			return
		}
		defer jsonFile.Close()

		b, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log.Println(err)
			return
		}

		err = json.Unmarshal(b, &cfg)
		if err != nil {
			log.Println(err)
			return
		}
	})
	return cfg
}
