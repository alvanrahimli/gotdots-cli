package utils

import "strings"

func SterilizeString(str string) string {
	forbiddenChars := []string{"*", ".", "\"", "/", "\\", "[", "]", ":", ";", "|", ",", "-"}
	for _, char := range forbiddenChars {
		str = strings.ReplaceAll(str, char, "_")
	}

	return str
}
