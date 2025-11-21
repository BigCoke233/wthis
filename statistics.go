package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type statistics struct {
	PkgName  string
	Name     string
	Type     string //formula or cask
	Desc     string
	Outdated bool

	InstalledAsDependency bool
	InstalledOnRequest    bool
	BinaryExists          bool
	LastAccessTime        time.Time

	InstalledVersion string
	Homepage         string
	License          string

	ReverseDependencies []string
}

const dataTTL = time.Hour / 2

func NewStatistics(formula *FormulaInfo, cask *CaskInfo, pkgName string, rvs []string) *statistics {
	var statObject statistics

	// === formula/cask specific === //

	if formula != nil {
		statObject.Name = formula.Name
		statObject.Type = "formula"
		statObject.Desc = formula.Desc
		statObject.Homepage = formula.Homepage
		statObject.License = formula.License
		statObject.Outdated = formula.Outdated

		// if installed
		if len(formula.Installed) != 0 {
			statObject.InstalledVersion = formula.Installed[0].Version
			statObject.InstalledAsDependency = formula.Installed[0].InstalledAsDependency
			statObject.InstalledOnRequest = formula.Installed[0].InstalledOnRequest
		}

		// get reverse dependencies
		statObject.ReverseDependencies = rvs
	} else if cask != nil {
		statObject.Name = cask.Token
		statObject.Type = "cask"
		statObject.Desc = cask.Desc
		statObject.Homepage = cask.Homepage
		statObject.InstalledVersion = cask.Installed
		statObject.Outdated = cask.Outdated
	}

	// === general === //

	statObject.PkgName = pkgName

	// get system atime
	var atime time.Time
	binPath := filepath.Join("/opt/homebrew/bin", pkgName)
	info, err := os.Stat(binPath)
	if err != nil {
		atime = time.Time{} // skip if no binary is found
		statObject.BinaryExists = false
	} else {
		atime = time.Unix(info.Sys().(*syscall.Stat_t).Atimespec.Unix())
		statObject.BinaryExists = true
	}
	statObject.LastAccessTime = atime

	return &statObject
}

func NewStatisticsFromCache(pkgName string) *statistics {
	raw, success := ReadCache("pkg/"+pkgName+".json", dataTTL)
	if success {
		var stat statistics
		json.Unmarshal(raw, &stat)
		return &stat
	}
	return nil
}

func (s *statistics) Cache() {
	data, err := json.Marshal(s)
	if err != nil {
		fmt.Println("Error marshaling statistics:", err)
		return
	}
	WriteCache("pkg/"+s.PkgName+".json", data)
}

func (s *statistics) Print() {
	if s == nil {
		fmt.Println("Package not found.")
		return
	}

	PrintHeader(s)
	PrintUserInteractionSummary(s)
	PrintMetadata(s)
	fmt.Println()
}
