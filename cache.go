package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
)

func CachePath(rel string) string {
	path := filepath.Join(xdg.CacheHome, "wthis", rel)
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	return path
}

func ReadCache(relPath string, ttl time.Duration) ([]byte, bool) {
	path := CachePath(relPath)

	info, err := os.Stat(path)
	if err != nil {
		return nil, false
	}
	if time.Since(info.ModTime()) > ttl {
		return nil, false
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}
	return b, true
}

func WriteCache(relPath string, data []byte) error {
	path := CachePath(relPath)
	return os.WriteFile(path, data, 0o644)
}
