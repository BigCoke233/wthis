package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/fatih/color"
)

const TTL = 3 * 24 * time.Hour // 3 days
const cacheName = "brew-checked"

func EnsureBrewAvailable() {
	// 1. check if cache file exists and is fresh
	data, _ := ReadCache(cacheName, TTL)
	t, _ := time.Parse(time.RFC3339, string(data))
	if time.Since(t) < TTL {
		return
	}

	// 2. otherwise, actually check
	if _, err := exec.LookPath("brew"); err != nil {
		// brew not found — no cache write
		color.Red("🍺 Homebrew not installed.")
		fmt.Println("This tool works only with Homebrew formulae and casks.")
		os.Exit(1)
	}

	// 3. write cache
	WriteCache(cacheName, []byte(time.Now().Format(time.RFC3339)))
}
