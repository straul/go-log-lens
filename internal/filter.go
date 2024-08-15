package internal

import (
	"regexp"
	"strings"
)

// FilterLogs filters logs based on a keyword.
func FilterLogs(logs []string, keywords []string) []string {
	var filteredLogs []string

	for _, line := range logs {
		for _, keyword := range keywords {
			if strings.Contains(line, keyword) {
				filteredLogs = append(filteredLogs, line)
			}
		}
	}

	return filteredLogs
}

// FilterLogsByRegex filters logs based on a regular expression pattern.
func FilterLogsByRegex(logs []string, pattern string) ([]string, error) {
	var filteredLogs []string

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	for _, line := range logs {
		if re.MatchString(line) {
			filteredLogs = append(filteredLogs, line)
		}
	}

	return filteredLogs, nil
}
