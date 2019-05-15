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

// GetConfig reads config file on startup and returns configuration
func GetConfig() *Config {
	doOnce.Do(func() {
		cfg = &Config{
			ServiceName:       "pic-resizer",
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
			log.Printf(`running config: 
			ServiceName: %v,
			Addr: %v,
			RootURL: %v,
			GetURL: %v,
			MultipartURL: %v,
			Base64URL: %v,
			ImagesURL: %v,
			ImagesPath: %v,
			ImagesPreviewSuff: %v`,
				cfg.ServiceName,
				cfg.Addr, cfg.RootURL,
				cfg.GetURL, cfg.MultipartURL,
				cfg.Base64URL, cfg.ImagesURL,
				cfg.ImagesPath, cfg.ImagesPreviewSuff)
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
