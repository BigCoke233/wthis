package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/fatih/color"
)

const TTL = 3 * 24 * time.Hour // 3 days

func EnsureBrewAvailable() {
	cache := filepath.Join(xdg.CacheHome, "wthis", "brew-checked")

	// 1. check if cache file exists and is fresh
	if info, err := os.Stat(cache); err == nil {
		if time.Since(info.ModTime()) < TTL {
			return // recently verified that brew exists
		}
	}

	// 2. otherwise, actually check
	if _, err := exec.LookPath("brew"); err != nil {
		// brew not found — no cache write
		color.Red("🍺 Homebrew not installed.")
		fmt.Println("This tool works only with Homebrew formulae and casks.")
		os.Exit(1)
	}

	// 3. write cache
	_ = os.MkdirAll(filepath.Dir(cache), 0o755)
	_ = os.WriteFile(cache, []byte{}, 0o644)
}
