package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// RegexFindTerm to find word in string
func RegexFindTerm(source, pattern string) string {
	regex := regexp.MustCompile(pattern)
	return regex.FindString(source)
}

// RegexFindBetweenTwoPattern do regex to get word inside two string
func RegexFindBetweenTwoPattern(before, after, sourceStr string) string {

	result := ""
	if !strings.Contains(sourceStr, before) {
		return result
	}

	beforeRegex := regexp.MustCompile(fmt.Sprintf(".*%v", before))
	afterRegex := regexp.MustCompile(fmt.Sprintf("%v.*", after))

	sourceStr = beforeRegex.ReplaceAllString(sourceStr, "")
	result = sourceStr
	if after != "" {
		result = afterRegex.ReplaceAllString(sourceStr, "")
	}

	return strings.TrimSpace(result)
}
