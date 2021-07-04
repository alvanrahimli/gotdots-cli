package main

import (
	"encoding/json"
	"fmt"
	"gotDots/models"
	"gotDots/utils"
	"io/fs"
	"os"
	"path"
	"path/filepath"
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
	fmt.Println("Help: dots <command> <options>")
	fmt.Println("  	new     <pack name> :	Creates new package with given name")
	fmt.Println("  	push    <pack name> :	Pushes package to registry (aws s3 for now)")
	fmt.Println("  	get     <pack name> :	Downloads package")
	fmt.Println("	install <pack name> :	Installs package to system")
}

func createManifest(packageName string, apps []GotDotsApp) models.Manifest {
	fmt.Print("Type version number (ex. 1.2.3): ")
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
			Version: app.GetVersion(),
		})
	}

	var visibility string
	choice, choiceErr := getYesNoChoice("Do you want to make package Public?")
	if choiceErr != nil {
		fmt.Println("Could not get choice")
		os.Exit(1)
	}

	if choice == models.YES {
		visibility = "Public"
	} else if choice == models.NO {
		visibility = "Private"
	}

	// TODO: Get userinfo
	var username = "USERNAME"
	var email = "EMAIL"

	return models.Manifest{
		Id:           fmt.Sprintf("org.gotdots.%s.%s", username, packageName),
		Name:         packageName,
		Version:      version.ToString(),
		Visibility:   visibility,
		IncludedApps: includedApps,
		Author: models.Author{
			Name:  username,
			Email: email,
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
	forbiddenChars := []string{"*", ".", "\"", "/", "\\", "[", "]", ":", ";", "|", ",", "-"}
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
		envErr := os.Setenv(lineSeperated[0], lineSeperated[1])
		if envErr != nil {
			return
		}
	}
}

// findPackageArchive returns archive's path for given package
func findPackageArchive(packName string) string {
	archiveFolder, folderErr := getArchivesFolder()
	if folderErr != nil {
		fmt.Printf("ERROR: %s\n", folderErr.Error())
		panic(folderErr)
	}

	var foundPack string
	walkErr := filepath.Walk(archiveFolder, func(path string, info fs.FileInfo, err error) error {
		// TODO: Refactor this to match whole package name
		if strings.Contains(path, packName) {
			foundPack = path
		}

		return nil
	})

	if walkErr != nil {
		fmt.Printf("ERROR: %s\n", walkErr.Error())
		panic(walkErr)
	}

	return foundPack
}

func readToken() string {
	archivesFolder, folderErr := getArchivesFolder()
	if folderErr != nil {
		panic(folderErr)
	}

	tokenStr, tokenErr := utils.ReadFromFile(path.Join(archivesFolder, "token"))
	if tokenErr != nil {
		panic(tokenErr)
	}

	return tokenStr
}

func handleError(err error, shouldExit bool) {
	fmt.Printf("ERROR: %s\n", err.Error())
	if shouldExit {
		os.Exit(1)
	}
}

func encodeIncludedApps(apps []models.IncludedApp) string {
	jsonApps, jsonErr := json.Marshal(apps)
	if jsonErr != nil {
		fmt.Println("Error occurred while marshalling included apps")
		fmt.Printf("ERROR: %s\n", jsonErr.Error())
		os.Exit(1)
	}

	return string(jsonApps)
}

func readManifestFile(fileAddress string) models.Manifest {
	manifestJson, readErr := os.ReadFile(fileAddress)
	if readErr != nil {
		handleError(readErr, true)
	}

	manifest := models.Manifest{}
	marshallErr := json.Unmarshal(manifestJson, &manifest)
	if marshallErr != nil {
		handleError(marshallErr, true)
	}

	return manifest
}

func getYesNoChoice(question string) (models.Choice, error) {
	fmt.Printf("%s (Y/n)", question)
	var choice string
	_, scanErr := fmt.Scanln(&choice)
	if scanErr != nil {
		return models.NO, scanErr
	}

	// Default is Y
	if choice == "" || choice == "Y" || choice == "y" {
		return models.YES, nil
	} else {
		return models.NO, nil
	}

}
