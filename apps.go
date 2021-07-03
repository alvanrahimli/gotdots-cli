package main

import (
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
	// Get absolute path of app's executable file
	appPath, err := exec.LookPath(appName)
	return appPath != "" && err == nil
}

// getRelativePath returns relative path of dotfile relative to ConfigRoot
// ex. Relative to "$HOME/.config/i3" for "$HOME/.config/i3/config" (returns "config")
func getRelativePath(app GotDotsApp, dotfile string) string {
	relativePath := strings.Replace(dotfile, app.GetConfigRoot(), "", -1)
	// This check makes sure path is not absolute
	if relativePath[0:1] == "/" {
		relativePath = relativePath[1:]
	}

	return relativePath
}
