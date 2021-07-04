package main

import (
	"fmt"
	"gotDots/utils"
	"os"
	"path"
)

func installPackage(packName string) {
	var packArchive string
	foundArchives := findPackageArchives(packName)
	if len(foundArchives) == 0 {
		fmt.Printf("Could not find package with name: %s\n", packName)
		os.Exit(1)
	} else if len(foundArchives) == 1 {
		packArchive = foundArchives[0]
	} else {
		fmt.Printf("Following packages found with name: %s\n", packName)
		utils.ListNames("   ", foundArchives)
		fmt.Print("Choose by entering number: ")
		var choice int
		_, scanErr := fmt.Scanln(&choice)
		if scanErr != nil {
			fmt.Println("Could not parse input")
			os.Exit(1)
		}

		packArchive = foundArchives[choice-1]
	}

	// Extract archive
	archiveFile, openErr := os.Open(packArchive)
	if openErr != nil {
		handleError(openErr, true)
	}

	packFolder := fmt.Sprintf("dots-pack-%s-*", packName)
	tempFolder, mkdirErr := os.MkdirTemp(os.TempDir(), packFolder)
	if mkdirErr != nil {
		handleError(mkdirErr, true)
	}

	untarErr := utils.Untar(tempFolder, archiveFile)
	if untarErr != nil {
		handleError(untarErr, true)
	}

	fmt.Println("Folder created at: " + tempFolder)

	supportedApps := getSupportedApps()
	// Read manifest file
	manifest := readManifestFile(path.Join(tempFolder, "manifest.json"))
	for _, app := range manifest.IncludedApps {
		isInstalled := isAppInstalled(app.Name)
		if !isInstalled {
			fmt.Printf("'%s' (%s) is not installed on your system. \n", app.Name, app.Version)
			// TODO: Handle if user wants to install app
			continue
		}

		installationErr := supportedApps[app.Name].InstallDotfiles(tempFolder, false)
		if installationErr != nil {
			handleError(installationErr, false)
		}
	}
}
