package helper

import "regexp"

func StripAllButNumbers(str string) string {
	// var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	var nonAlphanumericRegex = regexp.MustCompile(`[^0-9]+`)
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}
