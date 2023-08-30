package utils

import (
	"regexp"
)

func EscapeMD(s string) string {
	replacer := regexp.MustCompile(`[-_*.()\[\]~` + "`" + `#+=|{}!.]`)
	return replacer.ReplaceAllString(s, `\\$0`)
}
