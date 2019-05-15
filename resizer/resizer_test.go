package resizer

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/faroyam/pic-resizer/msg"
)

func TestStartService(t *testing.T) {
	StartService()
	Send(msg.Pic{ID: "testID", Data: make([]byte, 10)})
}

func TestResizerBadImage(t *testing.T) {
	err := resize(msg.Pic{ID: "testID", Data: make([]byte, 10)})
	if err == nil {
		t.Errorf("expected 'image: unknown format', got %v", err)
	}
}

func TestResizerGoodImage(t *testing.T) {
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

	err = os.Mkdir(cfg.ImagesPath, os.ModePerm)
	if err != nil {
		t.Errorf("cant create temp dir: %v", err)
	}

	defer func() {
		err = os.RemoveAll(cfg.ImagesPath)
		if err != nil {
			t.Errorf("cant delete temp dir: %v", err)
		}
	}()

	err = resize(msg.Pic{ID: "testID", Data: img})
	if err != nil {
		t.Errorf("%v", err)
	}
}
