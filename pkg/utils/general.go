package utils

// Ternary is a substitute for the ternary operator in other languages.
func Ternary[T any](predicate bool, onTrue T, onFalse T) T {
	if predicate {
		return onTrue
	} else {
		return onFalse
	}
}
