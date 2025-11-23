package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"log"
	"os"
)

var NoCache, DoClearCache bool

func main() {
	EnsureBrewAvailable()

	// declare CLI app
	cmd := &cli.Command{
		Name:  "wthis",
		Usage: "Fetch information of a Homebrew formulae or cask, like how it got installed, reverse dependencies, etc.",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "no-cache",
				Aliases:     []string{"nc"},
				Usage:       "Disable caching",
				Value:       false,
				Destination: &NoCache,
			},
			&cli.BoolFlag{
				Name:        "clear-cache",
				Aliases:     []string{"cc"},
				Usage:       "Clear cache",
				Value:       false,
				Destination: &DoClearCache,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if DoClearCache {
				ClearCache()
				fmt.Println("Cache cleared.")
				os.Exit(0)
			}

			// handle arguments
			if cmd.Args().Len() == 0 {
				printErrorAndExit("Please provide a formulae/cask name.")
			}
			pkgName := cmd.Args().First()

			// search and print
			showLoadingPrompt(fmt.Sprintf("What the heck is \"%s\"", pkgName))
			stat := fetchToCreateStatistics(pkgName)
			hideLoadingPrompt()
			stat.Print()

			// cache
			if !NoCache {
				stat.Cache()
			}

			return nil
		},
	}

	// run command and handle error
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func fetchToCreateStatistics(pkgName string) *statistics {
	// try reading from cache
	var stat *statistics
	if !NoCache {
		stat = NewStatisticsFromCache(pkgName)
	}
	// if cache not hit or explicitly requested no cache
	if NoCache || stat == nil {
		fmlChan := make(chan *FormulaInfo)
		caskChan := make(chan *CaskInfo)
		rvsChan := make(chan []string)
		// start 2 goroutines, fetching formula/cask info and uses
		go func() {
			formula, cask := GetBrewInfo(pkgName)
			fmlChan <- formula
			caskChan <- cask
		}()
		go func() {
			rvsChan <- GetBrewUses(pkgName)
		}()
		// create fine-grained data out of raw
		stat = NewStatistics(<-fmlChan, <-caskChan, pkgName, <-rvsChan)
	}
	return stat
}

func printErrorAndExit(msg string) {
	hideLoadingPrompt()
	color.Red("❌ %s", msg)
	os.Exit(1)
}
