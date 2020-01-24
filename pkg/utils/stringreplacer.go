package utils

import (
	"strings"
)

// StringReplacer represents the StringReplacer type
type StringReplacer string

// ReplaceAll finds and replaces all instances of 'replace' with 'with'
func (sr StringReplacer) ReplaceAll(replace, with string) StringReplacer {
	return StringReplacer(strings.ReplaceAll(string(sr), replace, with))
}

// String stringifies the StringReplacer
func (sr StringReplacer) String() string {
	return string(sr)
}
