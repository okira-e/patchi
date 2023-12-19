package utils

import "regexp"

// MaskString masks a string with asterisks.
func MaskString(str string) string {
	maskedStr := ""

	for i := 0; i < len(str); i++ {
		maskedStr += "*"
	}

	return maskedStr
}

// ExtractExpressions extracts expressions from a string.
func ExtractExpressions(str string, expr string) []string {
	// Extract the query from the template.
	regEx := regexp.MustCompile(expr)
	query := regEx.FindAllSubmatch([]byte(str), -1)

	result := make([]string, len(query))

	for i, q := range query {
		result[i] = string(q[1])
	}

	return result
}

// RemoveChar removes all occurrences of a character from a string.
func RemoveChar(str *string, char byte) {
	for i := 0; i < len(*str); i++ {
		if (*str)[i] == char {
			*str = (*str)[:i] + (*str)[i+1:]
			i--
		}
	}
}

// CapitalizeWord capitalizes the first letter of a word.
func CapitalizeWord(word string) string {
	if len(word) == 0 {
		return word
	}

	return string(word[0]-32) + word[1:]
}
