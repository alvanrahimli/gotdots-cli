package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

func CopyFile(sourceFile, destinationFile string) error {
	sourceFileStat, statErr := os.Stat(sourceFile)
	if statErr != nil {
		return statErr
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", sourceFile)
	}

	source, fileOpenErr := os.Open(sourceFile)
	if fileOpenErr != nil {
		return fileOpenErr
	}
	defer source.Close()

	destination, destErr := os.Create(destinationFile)
	if destErr != nil {
		return destErr
	}
	defer destination.Close()

	_, copyErr := io.Copy(destination, source)
	if copyErr != nil {
		return copyErr
	}

	return nil
}

func CopyFileToFolder(sourceFile, destinationFolder string) error {
	sourceFileStat, statErr := os.Stat(sourceFile)
	if statErr != nil {
		return statErr
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", sourceFile)
	}

	source, fileOpenErr := os.Open(sourceFile)
	if fileOpenErr != nil {
		return fileOpenErr
	}
	defer source.Close()

	_, fileName := path.Split(sourceFile)
	destination, destErr := os.Create(path.Join(destinationFolder, fileName))
	if destErr != nil {
		return destErr
	}
	defer destination.Close()

	_, copyErr := io.Copy(destination, source)
	if copyErr != nil {
		return copyErr
	}

	return nil
}

func WriteToFile(fileName, text string) error {
	writeErr := os.WriteFile(fileName, []byte(text), 0666)
	if writeErr != nil {
		return writeErr
	}

	return nil
}

func ReadFromFile(fileName string) (string, error) {
	fileContent, readErr := os.ReadFile(fileName)
	if readErr != nil {
		return "", readErr
	}

	return string(fileContent), nil
}

func DownloadFile(fileUrl, dest string) (string, error) {
	// Build fileName from fullPath
	fileURL, err := url.Parse(fileUrl)
	if err != nil {
		log.Fatal(err)
	}
	filePath := fileURL.Path
	segments := strings.Split(filePath, "/")
	fileName := segments[len(segments)-1]

	// Create blank file
	finalFile := path.Join(dest, fileName)
	file, err := os.Create(finalFile)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(fileUrl)
	if err != nil {
		return "", err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	size, err := io.Copy(file, resp.Body)

	defer file.Close()

	fmt.Printf("Downloaded a file %s with size %d\n", fileName, size)
	return finalFile, nil
}
