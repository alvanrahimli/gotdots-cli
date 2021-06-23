package utils

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type TarballFileType string

const (
	Dotfile TarballFileType = "dotfile"
	PackageFile TarballFileType = "pack_file"
)

//goland:noinspection GoUnhandledErrorResult
func CreateTarball(tempPackageFolder, tarballFilePath string) error {
	file, err := os.Create(tarballFilePath)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not create tarball file '%s', got error '%s'",
			tarballFilePath, err.Error()))
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()


	walkErr := filepath.Walk(tempPackageFolder, func(path string, info fs.FileInfo, err error) error {
		// Skip folders
		if info.IsDir() {
			return nil
		}

		relativePath := normalizePath(tempPackageFolder, path)
		tarWriterErr := addFileToTarWriter(path, relativePath, tarWriter)
		if tarWriterErr != nil {
			return tarWriterErr
		}

		return nil
	})

	if walkErr != nil {
		return walkErr
	}

	return nil
}

func addFileToTarWriter(filePath, relativePath string, tarWriter *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not open file '%s', got error '%s'",
			filePath, err.Error()))
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return errors.New(fmt.Sprintf("Could not get stat for file '%s', got error '%s'",
			filePath, err.Error()))
	}

	header := &tar.Header{
		Name:    relativePath,
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}

	err = tarWriter.WriteHeader(header)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not write header for file '%s', got error '%s'", filePath, err.Error()))
	}

	_, err = io.Copy(tarWriter, file)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not copy the file '%s' data to the tarball, got error '%s'", filePath, err.Error()))
	}

	return nil
}

func normalizePath(packFolder, dotfile string) string {
	dotfile = strings.Replace(dotfile, packFolder, "", -1)
	if dotfile[:1] == "/" {
		dotfile = dotfile[1:]
	}

	return dotfile
}