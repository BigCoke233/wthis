package main

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
)

// PrintInfo prints info for either a formula or a cask.
func PrintInfo(formula *FormulaInfo, cask *CaskInfo, pkgName string) {
	if formula == nil && cask == nil {
		fmt.Println("Package not found.")
		return
	}

	printHeader(formula, cask)
	printDescription(formula, cask)
	printVersionInfo(formula, cask)
	printMetadata(formula, cask)

	fmt.Println()
}

// ------------------- Helper functions -------------------

// printHeader prints the icon and package name
func printHeader(formula *FormulaInfo, cask *CaskInfo) {
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
	fmt.Printf(" - %s\n\n", desc)
}

// printDescription prints installation description for formulae
func printDescription(formula *FormulaInfo, cask *CaskInfo) {
	if formula != nil {
		if len(formula.Installed) == 0 {
			color.Blue("This package is not installed.\n")
		} else if formula.Installed[0].InstalledOnRequest {
			color.Blue("You installed this package by running `brew install`.\n")
		} else if formula.Installed[0].InstalledAsDependency {
			color.Blue("This package was installed automatically as a dependency.\n")
			// list reverse dependencies
			dependencies, err := GetReverseDependencies(formula.Name)
			if err != nil {
				fmt.Printf("Error getting reverse dependencies: %v\n", err)
			} else {
				fmt.Printf("Used by: %s\n", strings.Join(dependencies, ", "))
			}
		}
	} else if cask != nil {
		if cask.Installed == "" {
			color.Blue("This cask is not installed.\n")
		} else {
			color.Blue("This cask is installed.\n")
		}
	}

	fmt.Println()
}

// printVersionInfo prints the latest version and status
func printVersionInfo(formula *FormulaInfo, cask *CaskInfo) {
	if formula != nil {
		fmt.Printf("📦 Latest: %s", formula.Versions.Stable)
		if len(formula.Installed) == 0 {
			color.New(color.FgRed).Print(" [Not Installed]\n")
		} else if formula.Outdated {
			color.New(color.FgYellow).Print(" [Outdated]\n")
		} else {
			color.New(color.FgGreen).Print(" [Up to date]\n")
		}
	} else if cask != nil {
		fmt.Printf("📦 Installed version: %s", cask.Installed)
		if cask.Installed == "" {
			color.New(color.FgRed).Print(" [Not Installed]\n")
		} else if cask.Outdated {
			color.New(color.FgYellow).Print(" [Outdated]\n")
		} else {
			color.New(color.FgGreen).Print(" [Up to date]\n")
		}
	}
}

// printMetadata prints homepage, license, and other fields
func printMetadata(formula *FormulaInfo, cask *CaskInfo) {
	if formula != nil {
		fmt.Printf("🔗 Homepage: %s\n", formula.Homepage)
		fmt.Printf("📜 License: %s\n", formula.License)
	} else if cask != nil {
		fmt.Printf("🔗 Homepage: %s\n", cask.Homepage)
	}
}
