package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("=== Welcome to Got-Dots! ===")
	handleCommands(os.Args[1:])
}

func handleCommands(args []string) {
	if len(args) == 0 {
		printHelp()
		return
	}

	switch args[0] {
	case "new":
		if len(args) < 2 {
			fmt.Println("Please specify package name")
			return
		} else {
			createNewPackage(args[1])
		}
	case "install":
		if len(args) < 2 {
			fmt.Println("Please specify package name")
			return
		} else {
			installPackage(args[1])
		}
	case "help":
		printHelp()
	default:
		printHelp()
	}
}

func pushDotfiles(apps []GotDotsApp) {

}
