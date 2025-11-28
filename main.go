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
		},
		// Action: func(ctx context.Context, cmd *cli.Command) error {


		// 	return nil
		// },
		Commands: []*cli.Command{
			{
		       Name:    "info",
		       Aliases: []string{"i"},
		       Usage:   "Show info and reverse dependencies of a formula/cask.",
		       Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() == 0 {
						printErrorAndExit("Please provide a formulae/cask name.")
					}
					searchAndPrint(cmd.Args().First())

					return nil
		       },
		   },
			{
		       Name:    "list",
		       Aliases: []string{"l", "ls"},
		       Usage:   "Better brew list.",
		       Action: func(ctx context.Context, cmd *cli.Command) error {
					// TODO
					return nil
		       },
		   },
	        {
	            Name:    "clean",
	            Aliases: []string{"c", "cache"},
	            Usage:   "Delete all cache files.",
	            Action: func(ctx context.Context, cmd *cli.Command) error {
					// TODO error handling here
              		ClearCache()
					fmt.Println("Cache cleared.")
					return nil
	            },
	        },
		},
	}

	// run command and handle error
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func searchAndPrint(pkgName string) {
	// prompt
	showLoadingPrompt(fmt.Sprintf("What the heck is \"%s\"", pkgName))
	// try reading from cache
	var stat *statistics = NewStatisticsFromCache(pkgName)
	if NoCache || stat == nil {
		// cache not hit, then search
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
		stat = NewStatistics(<-fmlChan, <-caskChan, pkgName, <-rvsChan)
		// handle caching
		if !NoCache {
			stat.Cache()
		}
	}
	// print
	hideLoadingPrompt()
	stat.Print()
}

func printErrorAndExit(msg string) {
	hideLoadingPrompt()
	color.Red("❌ %s", msg)
	os.Exit(1)
}
