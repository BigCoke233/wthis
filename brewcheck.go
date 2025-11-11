package main

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"fmt"

	"github.com/fatih/color"
)

type BrewCache struct {
	Found     bool      `json:"found"`
	Timestamp time.Time `json:"timestamp"`
}

// Don't check Homebrew installation for a week
const cacheTTL = 7 * 24 * time.Hour

func EnsureBrewAvailable() {
	ok, err := isBrewAvailable()
	if err != nil {
		color.Red("⚠️  Error checking Homebrew: %v", err)
		os.Exit(1)
	}

	if !ok {
		color.Red("🍺 Homebrew not installed.")
		fmt.Println("This tool works only with Homebrew formulae and casks.")
		os.Exit(1)
	}
}

func isBrewAvailable() (bool, error) {
	cachePath, err := getCacheFile()
	if err != nil {
		return false, err
	}

	// Try to read cache
	if data, err := os.ReadFile(cachePath); err == nil {
		var cache BrewCache
		if json.Unmarshal(data, &cache) == nil {
			if time.Since(cache.Timestamp) < cacheTTL {
				return cache.Found, nil
			}
		}
	}

	// Perform actual check
	_, err = exec.LookPath("brew")
	found := err == nil

	// Save to cache
	_ = os.MkdirAll(filepath.Dir(cachePath), 0o755)
	cache := BrewCache{Found: found, Timestamp: time.Now()}
	data, _ := json.Marshal(cache)
	_ = os.WriteFile(cachePath, data, 0o644)

	return found, nil
}

func getCacheFile() (string, error) {
	base := os.Getenv("XDG_CACHE_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New("cannot determine user home")
		}
		base = filepath.Join(home, ".cache")
	}
	return filepath.Join(base, "wthis", "brew-check.json"), nil
}
