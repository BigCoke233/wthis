package main

import (
	"encoding/json"
	"github.com/fatih/color"
	"os/exec"
	"strings"
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
func GetBrewInfo(pkg string) (*FormulaInfo, *CaskInfo) {
	out, err := exec.Command("brew", "info", "--json=v2", pkg).Output()
	if err != nil {
		printErrorAndExit("Error fetching info. Please check for typo.")
	}

	var data struct {
		Formulae []FormulaInfo `json:"formulae"`
		Casks    []CaskInfo    `json:"casks"`
	}

	if err := json.Unmarshal(out, &data); err != nil {
		printErrorAndExit("Error parsing data.")
	}

	if len(data.Formulae) > 0 {
		return &data.Formulae[0], nil
	}
	if len(data.Casks) > 0 {
		return nil, &data.Casks[0]
	}

	return nil, nil // not found
}

func GetBrewUses(pkg string) []string {
	out, err := exec.Command("brew", "uses", "--installed", pkg).Output()
	if err != nil {
		color.Red("❌ Error fetching reverse dependencies.")
		return []string{}
	}

	// Each line is a formula name
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}
	}

	return lines
}
