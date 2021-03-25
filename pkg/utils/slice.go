package utils

import (
	"regexp"
)

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Filter(list []string, regex string) ([]string, error) {
	if regex == "" {
		return list, nil
	}

	filtered := make([]string, 0)
	matcher, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	for _, v := range list {
		if matcher.MatchString(v) {
			filtered = append(filtered, v)
		}
	}

	return filtered, nil
}
