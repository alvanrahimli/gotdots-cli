package utils

import (
	"fmt"
	"io"
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

func BatchCopyFiles(files []string, destinationFolder string) []string {
	var copiedFiles []string
	dotfileCount := len(files)
	for i, dotfile := range files {
		slashSplitted := strings.Split(dotfile, "/")
		fileName := slashSplitted[len(slashSplitted) - 1]
		copiedFileName := path.Join(destinationFolder, fileName)
		copyErr := CopyFile(dotfile, copiedFileName)
		if copyErr != nil {
			fmt.Printf("ERROR: %s\n", copyErr.Error())
		} else {
			copiedFiles = append(copiedFiles, copiedFileName)
			fmt.Printf("Copied %s (%d/%d)", fileName, i + 1, dotfileCount)
		}
	}

	return copiedFiles
}
