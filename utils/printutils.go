package utils

import (
	"fmt"
	"gotDots/models"
)

// ListNames prints names one per line with 1 based numbers
func ListNames(prefix string, names []string) {
	for i, name := range names {
		fmt.Printf("%s%d. %s\n", prefix, i+1, name)
	}
}

func GetYesNoChoice(question string, defaultChoice models.Choice) models.Choice {
	var alternateChoice models.Choice

	if defaultChoice == models.YES {
		alternateChoice = models.NO
		fmt.Printf("%s (Y/n) ", question)
	} else if defaultChoice == models.NO {
		alternateChoice = models.YES
		fmt.Printf("%s (y/N) ", question)
	}

	var choice string
	_, scanErr := fmt.Scanln(&choice)
	if scanErr != nil {
		return defaultChoice
	}

	// Default is Y
	if choice == "" || choice == "Y" || choice == "y" {
		return defaultChoice
	} else {
		return alternateChoice
	}
}
