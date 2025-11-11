package main

import (
	"encoding/json"
	"os/exec"
	"strings"
	"bytes"
	"github.com/fatih/color"
)

type FormulaInfo struct {
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Homepage string `json:"homepage"`
	License  string `json:"license"`
	Versions struct {
		Stable string `json:"stable"`
	} `json:"versions"`
	Installed []struct {
		Version               string `json:"version"`
		InstalledAsDependency bool   `json:"installed_as_dependency"`
		InstalledOnRequest    bool   `json:"installed_on_request"`
	} `json:"installed"`
	Outdated bool `json:"outdated"`
}

type CaskInfo struct {
	Token       string   `json:"token"`
	Name        []string `json:"name"`
	Desc        string   `json:"desc"`
	Homepage    string   `json:"homepage"`
	Installed   string   `json:"installed"`
	AutoUpdates bool     `json:"auto_updates"`
	Outdated    bool     `json:"outdated"`
	Version     string   `json:"version"`
}

// GetInfo fetches package info for both formulae and casks.
// It returns one of the two structs depending on the package type.
func GetInfo(pkg string) (*FormulaInfo, *CaskInfo) {
	out, err := exec.Command("brew", "info", "--json=v2", pkg).Output()
	if err != nil {
		color.Red("❌ Error fetching info. Please check for typo.")
		return nil, nil
	}

	var data struct {
		Formulae []FormulaInfo `json:"formulae"`
		Casks    []CaskInfo    `json:"casks"`
	}

	if err := json.Unmarshal(out, &data); err != nil {
		color.Red("❌ Error parsing data.")
		return nil, nil
	}

	if len(data.Formulae) > 0 {
		return &data.Formulae[0], nil
	}
	if len(data.Casks) > 0 {
		return nil, &data.Casks[0]
	}

	return nil, nil // not found
}

func GetReverseDependencies(pkg string) ([]string, error) {
	cmd := exec.Command("brew", "uses", "--installed", pkg)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Each line is a formula name
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}

	return lines, nil
}
