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
	// + TODO: Get package info from /packages/{packName}
	// + TODO: Download archive to $HOME/.dots-archives/
	// TODO: UnTar archive to /tmp/
	// TODO: Read included apps from manifest, check installation status
	// TODO: Install config files (Handlers should implement installation function)

	// Get package info
	formattedUrl := fmt.Sprintf("%s?name=%s", os.Getenv("GET_PACKAGE_URL"), packName)
	response, httpErr := http.Get(formattedUrl)
	if httpErr != nil {
		handleError(httpErr, true)
	}

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Error occured while fetching package info: %s", response.Status)
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

	fmt.Printf("Downloading package archive for: %s (%s)\n",
		packageInfo.Name, packVersionStr)

	// Download package archive
	archivesFolder, folderErr := getArchivesFolder()
	if folderErr != nil {
		handleError(folderErr, true)
	}

	packArchive, downloadErr := utils.DownloadFile(packageInfo.ArchiveUrl, archivesFolder)
	if downloadErr != nil {
		handleError(downloadErr, true)
	}

	archiveFile, openErr := os.Open(packArchive)
	if openErr != nil {
		handleError(openErr, true)
	}

	packFolder := path.Join("/tmp", packName)
	mkdirErr := os.Mkdir(packFolder, os.ModePerm)
	if mkdirErr != nil {
		handleError(mkdirErr, true)
	}

	untarErr := utils.Untar(packFolder, archiveFile)
	if untarErr != nil {
		handleError(untarErr, true)
	}

	// TODO: Continue this...
}
