package utils

import "strings"

func SplitStrings(str, sep string) []string {
	result := []string{}

	for _, part := range strings.Split(str, sep) {
		part = strings.TrimSpace(part)
		if len(part) == 0 {
			continue
		}

		result = append(result, part)
	}

	return result
}
