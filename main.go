package main

import (
    "fmt"
    "os"
    "os/exec"
    "context"
    "log"
    "github.com/fatih/color"
    "github.com/urfave/cli/v3"
)

func main() {
	// Check if Homebrew is installed
    _, brewInstallationErr := exec.LookPath("brew")
    if brewInstallationErr != nil {
        color.Red("🍺 Homebrew not installed.")
        fmt.Println("This tool is for looking up formulae/cask information about packages installed with Homebrew.")
        fmt.Println("If you don't have Homebrew, it's no use.")
        os.Exit(0)
    }

    // declare CLI app
    cmd := &cli.Command{
        Name:  "wthis",
        Usage: "Fetch information of a Homebrew formulae or cask, like how it got installed, reverse dependencies, etc.",
        Action: func(ctx context.Context, cmd *cli.Command) error {
        	if cmd.Args().Len() == 0 {
        		color.Red("❌ Please provide a formulae/cask name.")
          		fmt.Println("Usage: wthis <name>")
        		return nil
        	}
        	pkgName := cmd.Args().First()
         	formulae, cask := GetInfo(pkgName)
          	PrintInfo(formulae, cask, pkgName)
            return nil
        },
    }

    // run command and handle error
    if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}
