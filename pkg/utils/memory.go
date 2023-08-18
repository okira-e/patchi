package utils

// Copy returns a copy a value.
func Copy[T any](data *T) T {
	return *data
}
