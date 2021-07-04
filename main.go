package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("=== Welcome to Got-Dots! ===")
	loadEnvVariables()
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
			createNewPackage(sterilizeString(args[1]))
		}
	case "get":
		if len(args) < 2 {
			fmt.Println("Please specify package name")
			return
		} else {
			getPackage(args[1])
		}
	case "push":
		if len(args) < 2 {
			fmt.Println("Please specify package name")
			return
		} else {
			pushPackage(args[1])
		}
	case "update":
		if len(args) < 2 {
			fmt.Println("Please specify package name")
			return
		} else {
			updatePackage(args[1])
		}
	case "login":
		if len(args) < 1 {
			fmt.Println("Please specify package name")
			return
		} else {
			login()
		}
	case "help":
		printHelp()
	default:
		printHelp()
	}
}
