package main

import (
	"encoding/json"
	"fmt"
	"gotDots/models"
	"gotDots/utils"
	"os"
	"path"
)

func createNewPackage(packageName string) {
	foundApps := ScanForApps()
	if len(foundApps) == 0 {
		fmt.Println("No supported app found. \nPlease contact or consider contributing! <3")
		return
	}

	// Redirect to exclude apps dialogue
	foundApps = excludeApps(foundApps)

	manifest := createManifest(packageName, foundApps)

	// Make Tarball for package
	packageArchive, tarErr := createPackageArchive(manifest, foundApps)
	if tarErr != nil {
		fmt.Printf("ERROR: %s\n", tarErr.Error())
		fmt.Println("Could not create tarball")
		return
	}

	fmt.Println("Created new package at: " + packageArchive)
}

func createPackageArchive(manifest models.Manifest, apps []GotDotsApp) (string, error) {
	// Create temp folder for package
	folderNamePattern := fmt.Sprintf("gotdots-pack-%s-*", manifest.Name)
	tempPackageFolder, mkdirErr := os.MkdirTemp("/tmp", folderNamePattern)
	if mkdirErr != nil {
		fmt.Printf("Could not create temp folder. Error: %s\n", mkdirErr.Error())
		return "", mkdirErr
	}

	// Make package structure
	// Create dotfiles directory
	dotfilesFolderPath := path.Join(tempPackageFolder, "dotfiles")
	mkdirErr = os.Mkdir(dotfilesFolderPath, os.ModePerm)
	if mkdirErr != nil {
		fmt.Printf("Could not create dotfiles folder. Error: %s\n", mkdirErr.Error())
		return "", mkdirErr
	}

	var failedDotfiles []string
	// Copy files to appropriate folder in package folder
	for _, app := range apps {
		appDotfiles, dotfilesErr := app.GetExistingDotfiles()
		if dotfilesErr != nil {
			fmt.Printf("Could not get %s dotfiles. Error: %s", app.GetName(), dotfilesErr.Error())
			continue
		}

		// Create app specific folder
		appSpecificFolderName := path.Join(dotfilesFolderPath, app.GetName())
		mkdirErr = os.Mkdir(appSpecificFolderName, os.ModePerm)
		if mkdirErr != nil {
			fmt.Printf("Could not create dotfiles folder for app. Error: %s\n", mkdirErr.Error())
			continue
		}

		// Copy dotfiles for app to temp. package folder
		for _, dotfile := range appDotfiles {
			relativePath := getRelativePath(app, dotfile)
			copyErr := utils.CopyFile(dotfile, path.Join(appSpecificFolderName, relativePath))
			if copyErr != nil {
				fmt.Printf("Could not copy dotfile for app. Error: %s\n", copyErr.Error())
				failedDotfiles = append(failedDotfiles, dotfile)
				continue
			}
		}
	}

	// Print failed dotfiles
	if len(failedDotfiles) > 0 {
		fmt.Println("Could not include following dotfiles: ")
		for i, failedDotfile := range failedDotfiles {
			fmt.Printf("   %d. %s\n", i+1, failedDotfile)
		}
	}

	// Serialize manifest struct
	manifestJson, marshallErr := json.MarshalIndent(manifest, "", "  ")
	if marshallErr != nil {
		fmt.Println("Could not marshall manifest")
		return "", marshallErr
	}

	// Write manifest json to manifest.json file
	writeErr := os.WriteFile(path.Join(tempPackageFolder, "manifest.json"), manifestJson, os.ModePerm)
	if writeErr != nil {
		return "", writeErr
	}

	// Get archives folder
	archivesFolder, folderErr := getArchivesFolder()
	if folderErr != nil {
		fmt.Printf("ERROR: %s\n", folderErr.Error())
		return "", folderErr
	}

	// Create tar.gz file from temporary package folder
	// Looks like: gotdots-pack-%s-%s.tar.gz
	packTarName := fmt.Sprintf("gotdots-pack-%s-%s.tar.gz", manifest.Name, sterilizeString(manifest.Version))
	packTarFullName := path.Join(archivesFolder, packTarName)

	// Create tarball
	tarErr := utils.CreateTarball(tempPackageFolder, packTarFullName)
	if tarErr != nil {
		return "", tarErr
	}

	return packTarFullName, nil
}

func excludeApps(foundApps []GotDotsApp) []GotDotsApp {
	fmt.Println("Enter names to exclude from package (1 per line)")
	fmt.Println("Type 'x' to finish, 'l' to list remaining")

	appNames := getNames(foundApps)
	utils.ListNames("   ", appNames)

	// Excluding apps loop
	for {
		var input string
		fmt.Print("Enter name: ")
		_, err := fmt.Scan(&input)
		if err != nil || input == "x" {
			break
		} else if input == "l" {
			utils.ListNames("   ", appNames)
			continue
		}

		index := getIndexByName(foundApps, input)
		if index == -1 {
			fmt.Println("Could not find app. Try again")
			continue
		} else {
			foundApps = deleteByName(foundApps, input)
			appNames = getNames(foundApps)
			fmt.Printf("Excluded '%s'\n", input)
		}
	}

	return foundApps
}
