package utils

// MaskString masks a string with asterisks.
func MaskString(str string) string {
	maskedStr := ""

	for i := 0; i < len(str); i++ {
		maskedStr += "*"
	}

	return maskedStr
}
