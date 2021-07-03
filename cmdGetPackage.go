package main

import (
	"encoding/json"
	"fmt"
	"gotDots/models"
	"gotDots/utils"
	"io"
	"net/http"
	"os"
	"path"
)

func getPackage(packName string) {
	// Get package info from /packages/{packName}
	// Download archive to $HOME/.dots-archives/
	// UnTar archive to /tmp/
	// Read included apps from manifest, check installation status
	// Install config files (Handlers should implement installation function)

	// Get package info
	formattedUrl := fmt.Sprintf("%s?packageId=%s", os.Getenv("GET_PACKAGE_URL"), packName)
	response, httpErr := http.Get(formattedUrl)
	if httpErr != nil {
		handleError(httpErr, true)
	}

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Error occured while fetching package info: %s\n", response.Status)
		return
	}
	//goland:noinspection ALL
	defer response.Body.Close()

	body, bodyReadErr := io.ReadAll(response.Body)
	if bodyReadErr != nil {
		handleError(bodyReadErr, true)
	}

	var packageInfo models.PackageInfo
	unmarshallErr := json.Unmarshal(body, &packageInfo)
	if unmarshallErr != nil {
		handleError(unmarshallErr, true)
	}

	packVersionStr := fmt.Sprintf("%d.%d.%d",
		packageInfo.Version.Major, packageInfo.Version.Minor, packageInfo.Version.Patch)

	fmt.Printf("Downloading package archive for: %s (%s)\n", packageInfo.Name, packVersionStr)

	// Download package archive
	archivesFolder, folderErr := getArchivesFolder()
	if folderErr != nil {
		handleError(folderErr, true)
	}

	packArchive, downloadErr := utils.DownloadFile(packageInfo.ArchiveUrl, archivesFolder)
	if downloadErr != nil {
		handleError(downloadErr, true)
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
