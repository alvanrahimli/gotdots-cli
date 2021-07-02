package main

import (
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
		os.Setenv(lineSeperated[0], lineSeperated[1])
	}
}

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
	if shouldExit {
		os.Exit(1)
	}

	fmt.Printf("ERROR: %s\n", err.Error())
}

func encodeIncludedApps(apps []models.IncludedApp) string {
	var finalStr string
	// var innerStr []string
	for _, app := range apps {
		finalStr += fmt.Sprintf("[IncludedApps][Name]=%s&[IncludedApps][Version]=%s&", app.Name, app.Version)
		// innerStr = append(innerStr, fmt.Sprintf("%s:%s", app.Name, app.Version))
	}
	// finalStr = fmt.Sprintf("[%s]", strings.Join(innerStr, ","))

	return finalStr
}

func returnList(apps []models.IncludedApp) []string {
	var list []string
	for _, app := range apps {
		list = append(list, fmt.Sprintf("%s:%s", app.Name, app.Version))
	}

	return list
}
