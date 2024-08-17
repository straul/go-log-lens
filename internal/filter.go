package internal

import (
	"regexp"
	"strings"
	"time"
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

// FilterLogsByTime filters logs based on a time range.
func FilterLogsByTime(logs []string, startTime, endTime time.Time) []string {
	var filteredLogs []string
	timeRegex := regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`)

	for _, line := range logs {
		match := timeRegex.FindString(line)
		if match != "" {
			logTime, err := time.Parse(time.DateTime, match)
			if err == nil && logTime.After(startTime) && logTime.Before(endTime) {
				filteredLogs = append(filteredLogs, line)
			}
		}
	}

	return filteredLogs
}

// FilterLogsByLevel filters logs based on log levels.
func FilterLogsByLevel(logs []string, levels []string) []string {
	var filteredLogs []string

	for _, line := range logs {
		for _, level := range levels {
			if strings.Contains(line, level) {
				filteredLogs = append(filteredLogs, line)
				break
			}
		}
	}

	return filteredLogs
}

// FilterLine applies multiple filters to a single line and returns whether the line should be kept.
func FilterLine(line string, includeKeywords, excludeKeywords []string, levels []string, startTime, endTime *time.Time) bool {
	if startTime != nil && endTime != nil {
		timeRegex := regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`)
		match := timeRegex.FindString(line)
		if match != "" {
			logTime, err := time.Parse(time.DateTime, match)
			if err == nil && (logTime.Before(*startTime) || logTime.After(*endTime)) {
				return false
			}
		}
	}

	for _, keyword := range excludeKeywords {
		if strings.Contains(line, keyword) {
			return false
		}
	}

	if len(includeKeywords) > 0 && !containsAny(line, includeKeywords) {
		return false
	}

	if len(levels) > 0 && !containsAny(line, levels) {
		return false
	}

	return true
}

func containsAny(line string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(line, keyword) {
			return true
		}
	}
	return false
}
