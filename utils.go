package main

import (
	"fmt"
	"gotDots/models"
	"os"
	"path"
	"strings"
)

func getIndexByName(apps []GotDotsApp, name string) int {
	for i, app := range apps {
		if app.GetName() == name {
			return i
		}
	}

	return -1
}

func deleteByIndex(apps []GotDotsApp, index int) []GotDotsApp {
	apps[index] = apps[len(apps)-1]
	apps[len(apps)-1] = nil
	return apps[:len(apps)-1]
}

func deleteByName(apps []GotDotsApp, name string) []GotDotsApp {
	index := getIndexByName(apps, name)
	if index != -1 {
		return deleteByIndex(apps, index)
	}

	return apps
}

func getNames(apps []GotDotsApp) []string {
	var appNames []string
	for _, app := range apps {
		appNames = append(appNames, app.GetName())
	}

	return appNames
}

func printHelp() {
	fmt.Println("Help: ...")
}

func createManifest(packageName string, apps []GotDotsApp) models.Manifest {
	// TODO: Refactor this function

	fmt.Print("Type version number (ex. 1.0.4): ")
	version := models.PackageVersion{}
	_, scanErr := fmt.Scanf("%d.%d.%d", &version.Major, &version.Minor, &version.Patch)
	if scanErr != nil {
		fmt.Printf("ERROR: %s\n", scanErr.Error())
		version = models.PackageVersion{
			Major: 1,
			Minor: 0,
			Patch: 0,
		}
	}

	var includedApps []models.IncludedApp
	for _, app := range apps {
		includedApps = append(includedApps, models.IncludedApp{
			Name:    app.GetName(),
			Version: app.GetVersion().ToString(),
		})
	}

	return models.Manifest{
		Id:           "",
		Name:         packageName,
		Version:      version.ToString(),
		IncludedApps: includedApps,
		Author: models.Author{
			Name:  "PREDEFINED_VALUE",
			Email: "PREDEFINED_VALUE",
		},
	}
}

// getArchivesFolder returns $HOME/.dots-archives
func getArchivesFolder() (string, error) {
	// Get UserHomeDir
	userHomeDir, homeDirErr := os.UserHomeDir()
	if homeDirErr != nil {
		return "", homeDirErr
	}

	// Create archives folder
	archivesFolder := path.Join(userHomeDir, ".dots-archives")

	// Check if directory exists
	_, statErr := os.Stat(archivesFolder)
	if statErr != nil {
		// Create if directory does not exist
		mkdirErr := os.Mkdir(archivesFolder, os.ModePerm)
		if mkdirErr != nil {
			return "", mkdirErr
		}
	}

	return archivesFolder, nil
}

func sterilizeString(str string) string {
	forbiddenChars := []string{"*", ".", "\"", "/", "\\", "[", "]", ":", ";", "|", ","}
	for _, char := range forbiddenChars {
		str = strings.ReplaceAll(str, char, "_")
	}

	return str
}

func loadEnvVariables() {
	fileContent, readErr := os.ReadFile(".env")
	if readErr != nil {
		fmt.Printf("Could not read .env file. ERROR: %s\n", readErr.Error())
		return
	}

	variables := strings.Split(string(fileContent), "\n")
	for _, v := range variables {
		lineSeperated := strings.Split(v, "=")
		os.Setenv(lineSeperated[0], lineSeperated[1])
	}
}
