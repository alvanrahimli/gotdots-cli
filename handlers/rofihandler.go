package handlers

import (
	"fmt"
	"gotDots/utils"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
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

func (app RofiApp) GetVersion() string {
	output, outputErr := exec.Command("rofi", "-v").Output()
	if outputErr != nil {
		return "ERROR"
	}

	return strings.Split(string(output), " ")[1]
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
