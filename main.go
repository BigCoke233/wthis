package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

func main() {
	EnsureBrewAvailable()

	// declare CLI app
	cmd := &cli.Command{
		Name:  "wthis",
		Usage: "Fetch information of a Homebrew formulae or cask, like how it got installed, reverse dependencies, etc.",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// check args
			if cmd.Args().Len() == 0 {
				printErrorAndExit("Please provide a formulae/cask name.")
			}
			pkgName := cmd.Args().First()
			// prompt
			showLoadingPrompt(fmt.Sprintf("What the heck is \"%s\"", pkgName))
			// search, tidy up, and print
			fmlChan := make(chan *FormulaInfo)
			caskChan := make(chan *CaskInfo)
			rvsChan := make(chan []string)
			go func() {
				formula, cask := GetBrewInfo(pkgName)
				fmlChan <- formula
				caskChan <- cask
			}()
			go func() {
				rvsChan <- GetBrewUses(pkgName)
			}()
			stat := NewStatistics(<-fmlChan, <-caskChan, pkgName,<-rvsChan)
			hideLoadingPrompt()
			stat.Print()
			return nil
		},
	}

	// run command and handle error
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func printErrorAndExit(msg string) {
	hideLoadingPrompt()
	color.Red("❌ %s", msg)
	os.Exit(1)
}
