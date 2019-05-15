package config

import "testing"

func TestGetConfig(t *testing.T) {
	cfg := GetConfig()
	if cfg == nil {
		t.Errorf("cfg is nill")
	}
}
