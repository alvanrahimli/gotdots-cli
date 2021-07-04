package main

import (
	"fmt"
	"gotDots/utils"
	"os"
)

func updatePackage(packName string) {
	// Ignore these, as for now, user can only update their package
	//isInputCorrect, regexErr := regexp.MatchString("org\\.gotdots\\.[A-z0-9]{1,64}\\.[A-z0-9]{1,64}", packName)
	//if regexErr != nil {
	//	handleError(regexErr, true)
	//}
	//
	//if !isInputCorrect {
	//	fmt.Printf("Invalid package name format (%s) provided.\nIt should be\n" +
	//		"   org.gotdots.[username].[package name]\nTry again.\n", packName)
	//	os.Exit(1)
	//}

	author := readUserinfo()
	packageId := fmt.Sprintf("org.gotdots.%s.%s", author.Name, packName)
	foundArchives := findPackageArchives(packageId)

	var foundArchive string
	// Select archive from list
	if len(foundArchives) == 0 {
		fmt.Printf("Could not find package with name '%s'\n", packName)
		return
	} else if len(foundArchives) == 1 {
		foundArchive = foundArchives[0]
	} else if len(foundArchives) > 1 {
		fmt.Printf("Following packages found with name: %s\n", packName)
		utils.ListNames("   ", foundArchives)
		fmt.Print("Choose by entering number: ")
		var choice int
		_, scanErr := fmt.Scanln(&choice)
		if scanErr != nil {
			fmt.Println("Could not parse input")
			os.Exit(1)
		}

		foundArchive = foundArchives[choice-1]
	}

	manifest := utils.ReadManifestFromTar(foundArchive)
	fmt.Printf("Creating new version for package: %s\n", manifest.Id)
	createNewPackage(manifest.Name, manifest)
}
