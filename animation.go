package main

import (
	"fmt"
	"time"
	"strings"
	"github.com/fatih/color"
)

var loadingDone chan bool

func showLoadingPrompt(msg string) {
	if loadingDone != nil {
		return
	}

	loadingDone = make(chan bool)
	go func() {
		white := color.New(color.FgWhite)
		for {
			select {
			case <-loadingDone:
				fmt.Print("\r") // move cursor to line start
				fmt.Print(strings.Repeat(" ", len(msg)+4))
				fmt.Print("\r") // clear line
				return
			default:
				white.Printf("\r%s", msg)
				for i := range 4 {
					if i == 3 {
						white.Print("?")
					} else {
						white.Print(".")
					}
					time.Sleep(700 * time.Millisecond)
				}
			}
		}
	}()
}

// hideLoadingPrompt stops the animation cleanly.
func hideLoadingPrompt() {
	if loadingDone != nil {
		close(loadingDone)
		loadingDone = nil
		fmt.Print("\r\033[2K")
	}
}
