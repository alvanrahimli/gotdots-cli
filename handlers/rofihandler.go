package handlers

import (
	"fmt"
	"gotDots/models"
	"gotDots/utils"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

type RofiApp struct {
}

func (app RofiApp) GetPossibleDotfiles() []string {
	return []string{
		os.ExpandEnv("$HOME/.config/rofi/config.rasi"),
	}
}

func (app RofiApp) GetExistingDotfiles() ([]string, error) {
	var foundConfigs []string
	possibleConfigs := app.GetPossibleDotfiles()
	for _, configFile := range possibleConfigs {
		_, err := os.Stat(configFile)
		if !os.IsNotExist(err) {
			foundConfigs = append(foundConfigs, configFile)
		} else if os.IsPermission(err) {
			fmt.Printf("ERROR: Sufficient permission to read: %s\n", configFile)
		}
	}

	return foundConfigs, nil
}

func (app RofiApp) GetConfigRoot() string {
	return os.ExpandEnv("$HOME/.config/rofi")
}

func (app RofiApp) GetVersion() models.PackageVersion {
	// TODO: Implement this function
	return models.PackageVersion{
		Major: 1,
		Minor: 1,
		Patch: 1,
	}
}

func (app RofiApp) GetName() string {
	return "rofi"
}

func (app RofiApp) InstallDotfiles(packageFolder string, backup bool) error {
	fmt.Printf("Installing %s packages...\n", app.GetName())
	rofiDotfiles := path.Join(packageFolder, "dotfiles", app.GetName())
	walkErr := filepath.Walk(rofiDotfiles, func(path string, info fs.FileInfo, err error) error {

		//fmt.Println("Walking: " + path)
		if info.IsDir() {
			_, statErr := os.Stat(path)
			if os.IsNotExist(statErr) {
				mkdirErr := os.Mkdir(path, os.ModePerm)
				if mkdirErr != nil {
					return mkdirErr
				}
			}
		} else {
			copyErr := utils.CopyFileToFolder(path, app.GetConfigRoot())
			if copyErr != nil {
				return copyErr
			}

			return nil
		}

		return nil
	})
	if walkErr != nil {
		return walkErr
	}

	fmt.Printf("Finished installing for %s\n", app.GetName())
	return nil
}
