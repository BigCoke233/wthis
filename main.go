package main

import (
    "fmt"
    "os"
    "context"
    "log"
    "github.com/fatih/color"
    "github.com/urfave/cli/v3"
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
        		color.Red("❌ Please provide a formulae/cask name.")
          		fmt.Println("Usage: wthis <name>")
        		return nil
        	}
         	pkgName := cmd.Args().First()
         	// prompt
          	color.White("What the heck is \"%s\" ...?", pkgName)
         	// search and print
         	formulae, cask := GetBrewInfo(pkgName)
          	data := UnifyInfo(formulae, cask, pkgName)
          	PrintInfo(data)
            return nil
        },
    }

    // run command and handle error
    if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}
