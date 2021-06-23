package utils

import "fmt"

func ListNames(prefix string, names []string) {
	for i, name := range names {
		fmt.Printf("%s%d. %s\n", prefix, i + 1, name)
	}
}