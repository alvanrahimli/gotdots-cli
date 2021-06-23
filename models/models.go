package models

import "fmt"

type IncludedApp struct {
	Name	string
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
	Id           string         `json:"id"`
	Name         string         `json:"name"`
	Version      string 		`json:"version"`
	IncludedApps []IncludedApp 	`json:"includedApps"`
	Author       Author         `json:"author"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}