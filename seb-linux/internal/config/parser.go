package config

import (
	"os"

	"howett.net/plist"
)

func ParseSEBFile(path string) (*SEBConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &SEBConfig{KioskMode: true} // Kiosk por defecto
	if _, err := plist.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
