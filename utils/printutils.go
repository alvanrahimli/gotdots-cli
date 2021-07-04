package utils

import "fmt"

// ListNames prints names one per line with 1 based numbers
func ListNames(prefix string, names []string) {
	for i, name := range names {
		fmt.Printf("%s%d. %s\n", prefix, i+1, name)
	}
}
