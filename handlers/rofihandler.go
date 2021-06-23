package handlers

import (
	"fmt"
	"gotDots/models"
	"os"
)

type RofiApp struct {

}

func (app RofiApp) GetPossibleDotfiles() []string {
	return []string {
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
	return models.PackageVersion{
		Major: 1,
		Minor: 1,
		Patch: 1,
	}
}


func (app RofiApp) GetName() string {
	return "rofi"
}
