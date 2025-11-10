package main

import (
    "fmt"
    "os"
    "os/exec"
    "github.com/fatih/color"
)

func main() {
	// Check if Homebrew is installed
    _, brewInstallationErr := exec.LookPath("brew")
    if brewInstallationErr != nil {
        color.Red("🍺 Homebrew not installed.")
        fmt.Println("This tool is for looking up formulae/cask information about packages installed with Homebrew. If you don't have Homebrew, it's no use.")
        os.Exit(1)
    }

	args := os.Args
	// handle arguments
	var printUsage = func() {
		color.New(color.FgBlue).Print("Usage: ")
        fmt.Print("wtfis <package>\n")
	}
    if len(args) < 2 {
        printUsage()
        os.Exit(1)
    }
    if args[1] == "-h" || args[1] == "--help" {
    	printUsage()
        os.Exit(0)
    }
    pkg := args[1]

    // look up formulae
    out, brewLookupError := exec.Command("brew", "formulae", "--json=v2", pkg).Output()
    if brewLookupError != nil {
        fmt.Println("🍺 Homebrew Error:", string(out))
        os.Exit(1)
    }

    // parse formulae
    var formulae, cask, getformulaeErr = GetInfo(pkg)
    if getformulaeErr != nil {
        fmt.Println("Error:", getformulaeErr)
        os.Exit(1)
    }

    PrintInfo(formulae, cask, pkg)
}
