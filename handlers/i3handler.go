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

type I3WindowManager struct {
}

func (wm I3WindowManager) GetPossibleDotfiles() []string {
	return []string{
		os.ExpandEnv("$HOME/.config/i3/config"),
		os.ExpandEnv("$HOME/.i3/config"),
	}
}

func (wm I3WindowManager) GetExistingDotfiles() ([]string, error) {
	var foundConfigs []string
	possibleConfigs := wm.GetPossibleDotfiles()
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

func (wm I3WindowManager) GetConfigRoot() string {
	return os.ExpandEnv("$HOME/.config/i3")
}

func (wm I3WindowManager) GetVersion() models.PackageVersion {
	// TODO: Complete this to return real version
	return models.PackageVersion{
		Major: 1,
		Minor: 1,
		Patch: 1,
	}
}

func (wm I3WindowManager) GetName() string {
	return "i3"
}

func (wm I3WindowManager) InstallDotfiles(packageFolder string, backup bool) error {
	fmt.Printf("Installing %s packages...\n", wm.GetName())
	i3Dotfiles := path.Join(packageFolder, "dotfiles", wm.GetName())
	walkErr := filepath.Walk(i3Dotfiles, func(path string, info fs.FileInfo, err error) error {

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
			copyErr := utils.CopyFileToFolder(path, wm.GetConfigRoot())
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

	fmt.Printf("Finished installing for %s\n", wm.GetName())
	return nil
}
