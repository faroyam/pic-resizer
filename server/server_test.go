package server

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetImage(t *testing.T) {

	req, err := http.NewRequest("GET", "http://localhost"+cfg.Addr+cfg.RootURL+cfg.GetURL+"?url=http://thiscatdoesnotexist.com", nil)
	if err != nil {
		t.Errorf("creating request error: %v", err)
	}

	recorder := httptest.NewRecorder()
	getImageH(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected OK, got %v", recorder.Code)
	}
}

func TestGetImageBadURL(t *testing.T) {

	req, err := http.NewRequest("GET", "http://localhost"+cfg.Addr+cfg.RootURL+cfg.GetURL+"?url=http://thisURLdoesnotexist.com", nil)
	if err != nil {
		t.Errorf("creating request error: %v", err)
	}
	recorder := httptest.NewRecorder()
	getImageH(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected StatusBadRequest, got %v", recorder.Code)
	}
}

func TestBase64(t *testing.T) {

	resp, err := http.Get("http://thiscatdoesnotexist.com")
	if err != nil {
		t.Errorf("%v", err)
	}
	if resp != nil {
		defer resp.Body.Close()
	}

	img, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("%v", err)
	}

	imgBase64Str := base64.StdEncoding.EncodeToString(img)

	var jsonStr = []byte(fmt.Sprintf(`{"data":"%v"}`, imgBase64Str))
	req, err := http.NewRequest("POST", "http://localhost"+cfg.Addr+cfg.RootURL+cfg.Base64URL, bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Errorf("creating request error: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	base64H(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected OK, got %v", recorder.Code)
	}
}

func TestMultipartH(t *testing.T) {

	req, err := http.NewRequest("POST", "http://localhost"+cfg.Addr+cfg.RootURL+cfg.MultipartURL, nil)
	if err != nil {
		t.Errorf("creating request error: %v", err)
	}
	recorder := httptest.NewRecorder()
	multipartH(recorder, req)
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected OK, got %v", recorder.Code)
	}
}

func TestStartServe(t *testing.T) {

	err := make(chan error)
	go func() {
		err <- Serve()
	}()
	select {
	case e := <-err:
		t.Errorf("serve error: %v", e)
	case <-time.After(time.Second):
	}
}
