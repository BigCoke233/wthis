package main

import (
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type statistics struct {
	Name string;
	Type string;		//formula or cask
	Desc string;
	Outdated bool;

	InstalledAsDependency bool;
	InstalledOnRequest bool;
	LastAccessTime time.Time;

	InstalledVersion string;
	Homepage string;
	License string;

	ReverseDependencies []string;
}

func UnifyInfo(formula *FormulaInfo, cask *CaskInfo, pkgName string) (*statistics) {
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
		statObject.ReverseDependencies = GetBrewUses(pkgName)
	} else if cask != nil {
		statObject.Name = cask.Token
		statObject.Type = "cask"
		statObject.Desc = cask.Desc
		statObject.Homepage = cask.Homepage
		statObject.InstalledVersion = cask.Installed
		statObject.Outdated = cask.Outdated
	}

	// === general === //

	// get system atime
	binPath := filepath.Join("/opt/homebrew/bin", pkgName)
	info, err := os.Stat(binPath)
	if err != nil {
		return nil // skip if no binary is found
	}
	atime := info.Sys().(*syscall.Stat_t).Atimespec
	statObject.LastAccessTime = time.Unix(atime.Unix())

	return &statObject
}
