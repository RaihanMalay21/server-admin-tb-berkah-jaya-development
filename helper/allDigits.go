package helper

import "unicode"

func IsAllDigits(str string) bool {
	for _, char := range str {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}