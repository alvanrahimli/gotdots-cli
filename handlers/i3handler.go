package handlers

import (
	"fmt"
	"gotDots/models"
	"os"
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
