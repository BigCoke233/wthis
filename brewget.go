package main

import (
	"encoding/json"
	"os/exec"
	"strings"
	"time"

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

const brewDepCacheTTL = 3*24*time.Hour
const brewListCacheTTL = 24*time.Hour

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

func GetBrewList() []string {
	const cachepath = "brewlist.json"
	var list []string

	// Try reading from cache
	list_raw, cache_hit := ReadCache(cachepath, brewListCacheTTL)
	if cache_hit {
		err := json.Unmarshal(list_raw, &list)
		if err != nil {
			color.Red("❌ Error unmarshaling cache of brew list")
			return []string{}
		}
	} else {
		out, err := exec.
			Command("brew", "list", "--installed-on-request").
			Output()
		if err != nil {
			color.Red("❌ Error fetching brew list.")
			return []string{}
		}

		list = strings.Split(strings.TrimSpace(string(out)), "\n")
	}

	// Write cache
	defer func() {
		brewlist_json, err := json.Marshal(list)
		if err != nil {
			color.Red("❌ Error marshaling brew list data to write cache.")
		}
		WriteCache(cachepath, brewlist_json)
	}()

	return list
}

func GetBrewDeps() map[string][]string {
	depmap := map[string][]string{}

	// Try reading from cache
	const cachepath = "depmap.json"
	cache_data, cache_hit := ReadCache(cachepath, brewDepCacheTTL)

	if cache_hit {
		// if cache hit, decode
		err := json.Unmarshal(cache_data, &depmap)
		if err != nil {
			color.Red("❌ Error unmarshaling cache of dependency map.")
		}
	} else {
		// if not, try fetch
		out, err := exec.
			Command("brew", "deps", "--installed", "--direct").
			Output()
		if err != nil {
			color.Red("❌ Error fetching brew dependencies.")
			return map[string][]string{}
		}

		// parse raw data
		// a list of items like "package: depA depB depC..."
		dep_strs := strings.FieldsFunc(string(out), func(r rune) bool {
		    return r == '\n' || r == '\r'
		})

		for _, dep_str := range dep_strs {
			parts := strings.Split(dep_str, ":")
			pkg := parts[0]
			pkg_deps := parts[1]

			deps := strings.Split(pkg_deps, " ")
			filtered := []string{}
			// remove empty elements in slice
			for _, dep := range deps {
				if dep != "" {
					filtered = append(filtered, dep)
				}
			}
			depmap[pkg] = filtered
		}
	}

	// Writing cache
	defer func() {
		depmap_json, err := json.Marshal(depmap)
		if err != nil {
			color.Red("Error marshaling JSON when writing cache:", err)
		}
		WriteCache("depmap.json", depmap_json)
	}()

	return depmap
}
