package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func LoadAndConvert(path string) (App, error) {
	b, err := Load(path)
	if err != nil {
		return App{}, err
	}
	return Convert(b)
}

func Load(path string) ([]byte, error) {
	return os.ReadFile(filepath.Clean(path))
}

func Convert(b []byte) (App, error) {
	c := New()
	if err := json.Unmarshal(b, &c); err != nil {
		return c, fmt.Errorf("config.Convert: %v", err)
	}
	return c, nil
}
