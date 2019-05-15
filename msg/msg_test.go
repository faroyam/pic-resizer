package msg

import (
	"net/http/httptest"
	"testing"
)

func TestNewResponse(t *testing.T) {
	recorder := httptest.NewRecorder()
	err := NewResponse(recorder, "test", "original", "resized")
	if err != nil {
		t.Errorf("creating message error: %v", err)
	}
}
