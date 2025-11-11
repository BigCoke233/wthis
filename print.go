package main

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	"os"
	"path/filepath"
	"syscall"
	"time"
	"github.com/hako/durafmt"
)

// === Print Entry Point === //

func PrintInfo(formula *FormulaInfo, cask *CaskInfo, pkgName string) {
	if formula == nil && cask == nil {
		fmt.Println("Package not found.")
		return
	}

	printHeader(formula, cask)
	printUserInteractionSummary(formula, cask, pkgName)
	printMetadata(formula, cask)

	fmt.Println()
}

// === Main Sections === //

func printHeader(formula *FormulaInfo, cask *CaskInfo) {
	printTypeNameAndDesc(formula, cask)
	fmt.Println()
}

func printUserInteractionSummary(
	formula *FormulaInfo, cask *CaskInfo, pkgName string) {
	printInstallDesc(formula, cask)
	printRecentActivity(formula, cask, pkgName)
	fmt.Println()
}

func printMetadata(formula *FormulaInfo, cask *CaskInfo) {
	printVersionInfo(formula, cask)
	printHomeAndLicense(formula, cask)
}

// === #Section: Header# === //

func printTypeNameAndDesc(formula *FormulaInfo, cask *CaskInfo) {
	var icon, typ, name, desc string

	if formula != nil {
		icon = "🍺"
		typ = "(formula)"
		name = formula.Name
		desc = formula.Desc
	} else if cask != nil {
		icon = "☕️"
		typ = "(cask)"
		name = cask.Token
		desc = cask.Desc
	} else {
		name = "(unknown)"
	}

	fmt.Printf("%s ", icon)
	color.New(color.FgYellow).Printf("%s ", typ)
	color.New(color.Bold).Printf("%s", name)
	fmt.Printf(" - %s\n", desc)
}

// === #Section: User Interaction Summary# === //

// answers questions like "how it got here?"
func printInstallDesc(formula *FormulaInfo, cask *CaskInfo) {
	if cask != nil || formula == nil {
		return	// formulae only
	}

	if len(formula.Installed) == 0 {
		return
	} else if formula.Installed[0].InstalledOnRequest {
		color.Blue("► You installed this package by running `brew install`.\n")
	} else if formula.Installed[0].InstalledAsDependency {
		color.Blue("► This package was installed automatically as a dependency.\n")
		// list reverse dependencies
		dependencies, err := GetReverseDependencies(formula.Name)
		if err != nil {
			fmt.Printf("Error getting reverse dependencies: %v\n", err)
		} else {
			fmt.Printf("Used by: %s\n", strings.Join(dependencies, ", "))
		}
	}
}

// answers questions like "have I used this recently?"
func printRecentActivity(formula *FormulaInfo, cask *CaskInfo, pkgName string) {
	// skip if not installed
	if (formula != nil && len(formula.Installed) == 0) ||
		(cask != nil && cask.Installed == "") {
		return
	}

	// get last access time (atime)
	binPath := filepath.Join("/opt/homebrew/bin", pkgName)
	info, err := os.Stat(binPath)
	if err != nil {
		return	// skip if no binary is found
	}
	atime := info.Sys().(*syscall.Stat_t).Atimespec

	// print usage info
	var format string
	var humanReadableDuration string

	timestamp := time.Unix(atime.Unix())
	humanReadableDuration = durafmt.ParseShort(time.Since(timestamp)).String()
	format = func() string {
		if atime.Nano() == 0 {
			return "► You never used this."
		} else if time.Since(timestamp) < 24*7*time.Hour {
			return "► You used this in %s."
		} else {
			return "► You haven't used this for %s."
		}
	}()

	color.Blue(format, humanReadableDuration)
}

// === #Section: Metadata# === //

// prints badges like [Outdated] [Up to date] ...
// used by printMetadata
func printVersionInfo(formula *FormulaInfo, cask *CaskInfo) {
	var installedVersion string
	var installBadge string
	var installBadgeColor color.Attribute

	if formula != nil {
		if len(formula.Installed) != 0  {
			installedVersion = formula.Installed[0].Version
			if formula.Outdated {
				installBadge = "[Outdated]"
				installBadgeColor = color.FgYellow
			} else {
				installBadge = "[Up to date]"
				installBadgeColor = color.FgGreen
			}
		} else {
			installBadge = "[Not installed]"
			installBadgeColor = color.FgRed
		}
	} else if cask != nil {
		installedVersion = cask.Version
		if cask.Installed != "" {
			if cask.Outdated {
				installBadge = "[Outdated]"
				installBadgeColor = color.FgYellow
			} else {
				installBadge = "[Up to date]"
				installBadgeColor = color.FgGreen
			}
		} else {
			installBadge = "[Not installed]"
			installBadgeColor = color.FgRed
		}
	}

	fmt.Printf("📦 %s", installedVersion)
	color.New(installBadgeColor).Printf(" %s\n", installBadge)
}

func printHomeAndLicense(formula *FormulaInfo, cask *CaskInfo) {
	if formula != nil {
		fmt.Printf("🔗 %s\n", formula.Homepage)
		fmt.Printf("📜 %s\n", formula.License)
	} else if cask != nil {
		fmt.Printf("🔗 %s\n", cask.Homepage)
	}
}
