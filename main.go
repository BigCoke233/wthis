package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

// flags
var NoCache, DoClearCache bool
var ListAscend, ListDescend, ListNoDep bool

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
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				printErrorAndExit("Please provide a formulae/cask name.")
			}
			searchAndPrint(cmd.Args().First())

			return nil
		},
		Commands: []*cli.Command{
			{
		      	Name:    "list",
				Aliases: []string{"l", "ls"},
		   		Usage:   "Better brew list.",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "ascending",
						Aliases:     []string{"a"},
						Usage:       "Sort list ascendingly, based on dependency count.",
						Value:       false,
						Destination: &ListAscend,
					},
					&cli.BoolFlag{
						Name:        "descending",
						Aliases:     []string{"d"},
						Usage:       "Sort list descendingly, based on dependency count.",
						Value:       false,
						Destination: &ListDescend,
					},
					&cli.BoolFlag{
						Name:        "no-dependency",
						Aliases:     []string{"n", "nd"},
						Usage:       "Only show packages installed on request, with no dependency list.",
						Value:       false,
						Destination: &ListNoDep,
					},
				},
		       Action: func(ctx context.Context, cmd *cli.Command) error {
					depmap := GetBrewDeps()
					list := GetBrewList()

					// filter and get essential items
					var keys = []string{}
					for k := range depmap {
						if !slices.Contains(list, k) {
							delete(depmap, k)
						} else {
							keys = append(keys, k)
						}
					}

					// sort based on dependency count
					if ListAscend {
						sort.Slice(keys, func(i, j int) bool {
							return len(depmap[keys[i]]) > len(depmap[keys[j]])
						})
					} else if ListDescend {
						sort.Slice(keys, func(i, j int) bool {
							return len(depmap[keys[i]]) < len(depmap[keys[j]])
						})
					} else {
						sort.Strings(keys)
					}

					//
					for _, pkg := range keys {
						deps := depmap[pkg]
						count := len(deps)

						color.New(color.FgBlue).Printf("%s", pkg)
						if count > 1 {
							fmt.Printf(" (%d dependencies)\n", count)
						} else {
							fmt.Printf(" (%d dependency)\n", count)
						}

						if ListNoDep { continue }
						for _, dep := range deps {
							fmt.Printf("  ► %s\n", dep)
						}
					}
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
