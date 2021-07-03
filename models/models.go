package models

import (
	"fmt"
)

type IncludedApp struct {
	Name    string
	Version string
}

type PackageVersion struct {
	Major int
	Minor int
	Patch int
}

func (p PackageVersion) ToString() string {
	return fmt.Sprintf("%d.%d.%d", p.Major, p.Minor, p.Patch)
}

type Manifest struct {
	Id           string        `json:"id"`
	Name         string        `json:"name"`
	Version      string        `json:"version"`
	Visibility   string        `json:"visibility"`
	IncludedApps []IncludedApp `json:"includedApps"`
	Author       Author        `json:"author"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type PackageInfo struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	ArchiveUrl string `json:"archiveUrl"`
	Version    struct {
		Major int `json:"major"`
		Minor int `json:"minor"`
		Patch int `json:"patch"`
	} `json:"version"`
	DownloadCount int `json:"downloadCount"`
	Rating        int `json:"rating"`
	IncludedApps  []struct {
		AppName string
		Version struct {
			Major int `json:"major"`
			Minor int `json:"minor"`
			Patch int `json:"patch"`
		} `json:"version"`
	} `json:"includedApps"`
}

type Choice string

const (
	YES Choice = "yes"
	NO  Choice = "no"
)
