package main

import (
	"fmt"
	"github.com/fatih/color"
	"strings"
	"time"
	"github.com/hako/durafmt"
)

// === Print Entry Point === //

func PrintInfo(data *statistics) {
	if data == nil {
		fmt.Println("Package not found.")
		return
	}

	printHeader(data)
	printUserInteractionSummary(data)
	printMetadata(data)

	fmt.Println()
}

// === Main Sections === //

func printHeader(data *statistics) {
	printTypeNameAndDesc(data.Type, data.Name, data.Desc)
	fmt.Println()
}

func printUserInteractionSummary(data *statistics) {
	printInstallDesc(
		data.Type == "formula",
		data.InstalledAsDependency,
		data.InstalledOnRequest,
		data.InstalledVersion,
		data.ReverseDependencies)
	printRecentActivity(
		data.InstalledVersion,
	 	data.LastAccessTime)
	fmt.Println()
}

func printMetadata(data *statistics) {
	printVersionInfo(
		data.Type,
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
	color.New(color.FgYellow).Printf("%s ", typ)
	color.New(color.Bold).Printf("%s", name)
	fmt.Printf(" - %s\n", desc)
}

// === #Section: User Interaction Summary# === //

// answers questions like "how it got here?"
func printInstallDesc(
	isFormulae bool,
	asDependency bool,
	onRequest bool,
	version string,
	reverseDependencies []string,
) {
	// formulae only
	if !isFormulae {
		return
	}

	if version == "" {
		return
	} else if onRequest {
		color.Blue("► You installed this package by running `brew install`.\n")
	} else if asDependency {
		color.Blue("► This package was installed automatically as a dependency.\n")
		fmt.Printf("Used by: %s\n",
			strings.Join(reverseDependencies, ", "))
	}
}

// answers questions like "have I used this recently?"
func printRecentActivity(installedVersion string, atime time.Time) {
	// skip if not installed
	if installedVersion == "" {
		return
	}

	// print usage info
	humanReadableDuration := durafmt.ParseShort(time.Since(atime)).String()
	if atime.Nanosecond() == 0 {
		color.Blue("► You never used this.")
	} else if time.Since(atime) < 24*7*time.Hour {
		color.Blue("► You used this in %s.", humanReadableDuration)
	} else {
		color.Blue("► You haven't used this for %s.", humanReadableDuration)
	}
}

// === #Section: Metadata# === //

// prints badges like [Outdated] [Up to date] ...
// used by printMetadata
func printVersionInfo(typ string, installedVersion string, outdated bool) {
	var installBadge string
	var installBadgeColor color.Attribute

	if installedVersion != "" {
		if outdated {
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

	fmt.Printf("📦 %s", installedVersion)
	color.New(installBadgeColor).Printf(" %s\n", installBadge)
}

func printHomeAndLicense(homepage string, license string) {
	if homepage != "" {
		fmt.Printf("🔗 %s\n", homepage)
	}
	if license != "" {
		fmt.Printf("📜 %s\n", license)
	}
}
