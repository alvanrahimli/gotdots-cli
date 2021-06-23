package utils

import (
	"fmt"
	"io"
	"os"
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
