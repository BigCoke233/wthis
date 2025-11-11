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
			return
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
	}

	fmt.Println()
}

// printVersionInfo prints the latest version and status
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

// printMetadata prints homepage, license, and other fields
func printMetadata(formula *FormulaInfo, cask *CaskInfo) {
	if formula != nil {
		fmt.Printf("🔗 %s\n", formula.Homepage)
		fmt.Printf("📜 %s\n", formula.License)
	} else if cask != nil {
		fmt.Printf("🔗 %s\n", cask.Homepage)
	}
}
