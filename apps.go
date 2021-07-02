package main

import (
	"bytes"
	"fmt"
	"gotDots/handlers"
	"gotDots/models"
	"os/exec"
	"strings"
)

func getSupportedApps() map[string]GotDotsApp {
	return map[string]GotDotsApp{
		"i3":   handlers.I3WindowManager{},
		"rofi": handlers.RofiApp{},
	}
}

type GotDotsApp interface {
	// GetPossibleDotfiles returns absolute paths of possible dotfiles, is taken from documentation
	GetPossibleDotfiles() []string
	// GetExistingDotfiles returns absolute paths of existing dotfiles for app
	GetExistingDotfiles() ([]string, error)
	// GetConfigRoot returns root directory where dotfile(s) are located
	GetConfigRoot() string
	// GetVersion returns version number of app
	GetVersion() models.PackageVersion
	// GetName returns name of app's executable file
	GetName() string
	// InstallDotfiles backups old files and installs new files
	InstallDotfiles(packageFolder string, backup bool) error
}

func ScanForApps() []GotDotsApp {
	var foundApps []GotDotsApp

	for appName, instance := range getSupportedApps() {
		exists := isAppInstalled(appName)
		if exists {
			foundApps = append(foundApps, instance)
		}
	}

	return foundApps
}

func isAppInstalled(appName string) bool {
	// TODO: Change this and use os.Stat(...)

	cmd := exec.Command("whereis", appName)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error occured: %s\n", err.Error())
		return false
	}

	output := out.String()
	colonSplitted := strings.Split(output, ":")

	return len(colonSplitted[1]) > 1
}

// getRelativePath returns relative path of dotfile relative to ConfigRoot
// ex. Relative to "$HOME/.config/i3" for "$HOME/.config/i3/config" (returns "config")
func getRelativePath(app GotDotsApp, dotfile string) string {
	relativePath := strings.Replace(dotfile, app.GetConfigRoot(), "", -1)
	// Make sure path is not absolute
	if relativePath[0:1] == "/" {
		relativePath = relativePath[1:]
	}

	return relativePath
}
