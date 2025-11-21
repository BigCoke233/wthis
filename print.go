package main

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	"time"
	"github.com/hako/durafmt"
)

// === Main Sections === //

func PrintHeader(data *statistics) {
	printTypeNameAndDesc(data.Type, data.Name, data.Desc)
	fmt.Println()
}

func PrintUserInteractionSummary(data *statistics) {
	a := printInstallDesc(
		data.Type == "formula",
		data.InstalledAsDependency,
		data.InstalledOnRequest,
		data.InstalledVersion,
		data.ReverseDependencies)
	b := printRecentActivity(
		data.InstalledVersion,
	 	data.LastAccessTime,
		data.BinaryExists)

	if a || b {	// add empty line only if any of above prints info
		fmt.Println()
	}
}

func PrintMetadata(data *statistics) {
	printVersionInfo(
		data.InstalledVersion,
		data.Outdated)
	printHomeAndLicense(data.Homepage, data.License)
}

// === #Section: Header# === //

func printTypeNameAndDesc(typ string, name string, desc string) {
	var icon string
	switch typ {
		case "formula":
			icon = "🍺"
		case "cask":
			icon = "☕️"
	}

	fmt.Printf("%s ", icon)
	color.New(color.FgYellow).Printf("(%s) ", typ)
	color.New(color.Bold).Printf("%s", name)
	fmt.Printf(" - %s\n", desc)
}

// === #Section: User Interaction Summary# === //

// answers questions like "how it got here?"
func printInstallDesc(
	isFormula bool,
	asDependency bool,
	onRequest bool,
	version string,
	reverseDependencies []string,
) (printed bool) {
	// formula only
	if !isFormula {
		return false
	}

	if version == "" {
		return false
	} else if onRequest {
		color.Blue("► You installed this package by running `brew install`.\n")
	} else if asDependency {
		color.Blue("► This package was installed automatically as a dependency.\n")
		fmt.Printf("Used by: %s\n",
			strings.Join(reverseDependencies, ", "))
	}
	return true
}

// answers questions like "have I used this recently?"
func printRecentActivity(
	installed string,
	atime time.Time,
	binary bool) (printed bool) {

	// skip if not installed or no related binary
	if installed == "" || !binary {
		return false
	}

	// print usage info
	duratext := durafmt.ParseShort(time.Since(atime)).String()
	if atime.Nanosecond() == 0 {
		color.Blue("► You never used this.")
	} else if time.Since(atime) < 24*7*time.Hour {
		color.Blue("► You used this in %s.", duratext)
	} else {
		color.Blue("► You haven't used this for %s.", duratext)
	}
	return true
}

// === #Section: Metadata# === //

// prints badges like [Outdated] [Up to date] ...
// used by printMetadata
func printVersionInfo(installed string, outdated bool) {
	if installed != "" {
		var installBadge string
		var installBadgeColor color.Attribute
		if outdated {
			installBadge = "[Outdated]"
			installBadgeColor = color.FgYellow
		} else {
			installBadge = "[Up to date]"
			installBadgeColor = color.FgGreen
		}
		fmt.Printf("📦 %s", installed)
		color.New(installBadgeColor).Printf(" %s\n", installBadge)
	} else {
		color.Red("📦 [Not installed]")
	}

}

func printHomeAndLicense(homepage string, license string) {
	if homepage != "" {
		fmt.Printf("🔗 %s\n", homepage)
	}
	if license != "" {
		fmt.Printf("📜 %s\n", license)
	}
}
