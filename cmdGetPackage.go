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
		fmt.Printf("ERROR: %s\n", httpErr.Error())
		panic(httpErr)
	}

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Error occured while fetching package info: %s", response.Status)
		return
	}
	//goland:noinspection ALL
	defer response.Body.Close()

	body, bodyReadErr := io.ReadAll(response.Body)
	if bodyReadErr != nil {
		fmt.Printf("ERROR: %s\n", bodyReadErr.Error())
		// panic(bodyReadErr)
		return
	}

	var packageInfo models.PackageInfo
	unmarshallErr := json.Unmarshal(body, &packageInfo)
	if unmarshallErr != nil {
		fmt.Printf("ERROR: %s\n", unmarshallErr.Error())
		panic(unmarshallErr)
		// return
	}

	packVersionStr := fmt.Sprintf("%d.%d.%d",
		packageInfo.Version.Major, packageInfo.Version.Minor, packageInfo.Version.Patch)

	fmt.Printf("Downloading package archive for: %s (%s)\n",
		packageInfo.Name, packVersionStr)

	// Download package archive
	archivesFolder, folderErr := getArchivesFolder()
	if folderErr != nil {
		panic(folderErr)
	}

	packArchive, downloadErr := utils.DownloadFile(packageInfo.ArchiveUrl, archivesFolder)
	if downloadErr != nil {
		fmt.Printf("ERROR: %s\n", downloadErr.Error())
		panic(downloadErr)
	}

	archiveFile, openErr := os.Open(packArchive)
	if openErr != nil {
		panic(openErr)
	}

	packFolder := path.Join("/tmp", packName)
	mkdirErr := os.Mkdir(packFolder, os.ModePerm)
	if mkdirErr != nil {
		return
	}

	err := utils.Untar(packFolder, archiveFile)
	if err != nil {
		panic(err)
		// return
	}
}
