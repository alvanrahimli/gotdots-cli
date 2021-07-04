package main

import (
	"encoding/json"
	"fmt"
	"gotDots/models"
	"gotDots/utils"
	"io"
	"net/http"
	"os"
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

	_, downloadErr := utils.DownloadFile(packageInfo.ArchiveUrl, archivesFolder)
	if downloadErr != nil {
		handleError(downloadErr, true)
	}

	choice := utils.GetYesNoChoice("Do you want to install package now?", models.YES)

	if choice == models.YES {
		installPackage(packName)
	} else if choice == models.NO {
		fmt.Printf("You can install package by typing\n   dots install %s\nlater\n", packName)
	}
}
