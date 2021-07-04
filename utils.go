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
	choice := utils.GetYesNoChoice("Do you want to make package Public?", models.YES)

	if choice == models.YES {
		visibility = "Public"
	} else if choice == models.NO {
		visibility = "Private"
	}

	// TODO: Get userinfo
	author := models.Author{
		Name:  "ERROR",
		Email: "ERROR",
	}
	archivesFolder, _ := getArchivesFolder()
	authorStr, readErr := utils.ReadFromFile(path.Join(archivesFolder, ".userinfo"))
	if readErr != nil {
		handleError(readErr, false)
	}

	jsonErr := json.Unmarshal([]byte(authorStr), &author)
	if jsonErr != nil {
		handleError(jsonErr, false)
	}

	return models.Manifest{
		Id:           fmt.Sprintf("org.gotdots.%s.%s", author.Name, packageName),
		Name:         packageName,
		Version:      version.ToString(),
		Visibility:   visibility,
		IncludedApps: includedApps,
		Author: models.Author{
			Name:  author.Name,
			Email: author.Email,
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

// findPackageArchives returns archive's path for given package
// If multiple matches found, it asks user to choose from list
func findPackageArchives(packName string) []string {
	archiveFolder, folderErr := getArchivesFolder()
	if folderErr != nil {
		fmt.Printf("ERROR: %s\n", folderErr.Error())
		panic(folderErr)
	}

	var foundPacks []string
	walkErr := filepath.Walk(archiveFolder, func(folderMemberPath string, info fs.FileInfo, err error) error {
		// if it is archive file
		if !info.IsDir() && strings.Contains(folderMemberPath, "tar.gz") {
			manifest := utils.ReadManifestFromTar(folderMemberPath)
			if manifest.Id == packName {
				foundPacks = append(foundPacks, folderMemberPath)
			}
		}

		return nil
	})

	if walkErr != nil {
		fmt.Printf("ERROR: %s\n", walkErr.Error())
		panic(walkErr)
	}

	return foundPacks
}

func readToken() string {
	archivesFolder, folderErr := getArchivesFolder()
	if folderErr != nil {
		panic(folderErr)
	}

	tokenStr, tokenErr := utils.ReadFromFile(path.Join(archivesFolder, ".token"))
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
