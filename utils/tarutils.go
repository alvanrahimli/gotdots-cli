package utils

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"gotDots/models"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type TarballFileType string

//goland:noinspection GoUnhandledErrorResult
func CreateTarball(sourceFolder, tarballFilePath string) error {
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

	walkErr := filepath.Walk(sourceFolder, func(path string, info fs.FileInfo, err error) error {
		// Skip folders
		if info.IsDir() {
			return nil
		}

		// Normalize path, so it is not absolute any more
		relativePath := normalizePath(sourceFolder, path)

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

func Untar(dst string, r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {
		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)
		// fmt.Println("NOW TARGET: " + target)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.Create(target)
			if err != nil {
				fmt.Println("Could not create " + target + ". Creating parent directories")
				folders := strings.Split(target, "/")
				mkdirAllErr := os.MkdirAll(strings.Join(folders[:len(folders)-1], "/"), os.ModePerm)
				if mkdirAllErr != nil {
					fmt.Println("Could not create parent directories. Skipping this file")
					fmt.Printf("ERROR: %s\n", mkdirAllErr.Error())
					continue
				}
				f, _ = os.Create(target)
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

func ReadManifestFromTar(tarFileName string) models.Manifest {
	archiveFile, archiveOpenErr := os.Open(tarFileName)
	if archiveOpenErr != nil {
		fmt.Println("Error occurred while opening archive")
		fmt.Printf("ERROR: %s\n", archiveOpenErr.Error())
		os.Exit(1)
	}

	gzr, err := gzip.NewReader(archiveFile)
	if err != nil {
		fmt.Println("Error occurred while opening gzip reader")
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	var manifest models.Manifest

	for {
		header, err := tr.Next()
		if err != nil {
			fmt.Println("Error occurred while getting next header from tar")
			fmt.Printf("ERROR: %s\n", err.Error())
			continue
		}

		if header.Name == "manifest.json" {
			tempManifestFile, err := os.CreateTemp("/tmp", "")
			defer tempManifestFile.Close()

			if err != nil {
				fmt.Println("Error occurred while creating temp manifest file")
				fmt.Printf("ERROR: %s\n", err.Error())
				os.Exit(1)
			}

			targetFile, fileErr := os.OpenFile(tempManifestFile.Name(), os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if fileErr != nil {
				fmt.Println("Error occurred while opening/creating file while untar process")
				fmt.Printf("ERROR: %s\n", fileErr.Error())
				os.Exit(1)
			}

			if _, err := io.Copy(targetFile, tr); err != nil {
				fmt.Println("Error occurred copying file while untar process")
				fmt.Printf("ERROR: %s\n", err.Error())
			}

			targetFile.Close()

			fileReader, _ := os.ReadFile(tempManifestFile.Name())
			unmarshallErr := json.Unmarshal(fileReader, &manifest)
			if unmarshallErr != nil {
				fmt.Println("Error occurred while unmarshalling manifest file")
				fmt.Printf("ERROR: %s\n", unmarshallErr.Error())
				os.Exit(1)
			}

			return manifest
		}
	}

}
