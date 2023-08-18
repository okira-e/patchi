package collections

import "errors"

type Vec[T comparable] struct {
	data []T
}

// NewVec returns a new vector with the given data.
func NewVec[T comparable](data ...T) *Vec[T] {
	v := Vec[T]{}

	for _, val := range data {
		v.data = append(v.data, val)
	}

	return &v
}

// Len returns the length of the vector.
func (v *Vec[T]) Len() int {
	return len(v.data)
}

// Cap returns the capacity of the vector.
func (v *Vec[T]) Cap() int {
	return cap(v.data)
}

// Get returns the value at the given index.
func (v *Vec[T]) Get(index int) T {
	return v.data[index]
}

// Set sets the value at the given index.
func (v *Vec[T]) Set(index int, val T) error {
	if index >= len(v.data) {
		return errors.New("index out of range")
	}

	v.data[index] = val

	return nil
}

// Push appends the value to the end of the vector.
func (v *Vec[T]) Push(val T) {
	v.data = append(v.data, val)
}

// Remove removes the value from the vector.
func (v *Vec[T]) Remove(val T) {
	for i, _v := range v.data {
		if _v == val {
			v.data = append(v.data[:i], v.data[i+1:]...)
			break
		}
	}
}

// Pop removes the last value from the vector and returns it.
func (v *Vec[T]) Pop() T {
	last := v.data[len(v.data)-1]
	v.data = v.data[:len(v.data)-1]

	return last
}

func (v *Vec[T]) Range(start, end int) *Vec[T] {
	return NewVec[T](v.data[start:end]...)
}

// Shift removes the first value from the vector and returns it.
func (v *Vec[T]) Shift() T {
	first := v.data[0]
	v.data = v.data[1:]

	return first
}

// Unshift prepends the value to the beginning of the vector.
func (v *Vec[T]) Unshift(val T) {
	v.data = append([]T{val}, v.data...)
}

// ForEach iterates over the vector and calls the callback function for each value.
func (v *Vec[T]) ForEach(cb func(int, T)) {
	for i, val := range v.data {
		cb(i, val)
	}
}

// Map iterates over the vector and calls the callback function for each value.
// The callback function should return the new value of comparable type based on the current value.
// It returns a new vector with the new values.
func (v *Vec[T]) Map(cb func(int, T) T) *Vec[T] {
	newData := make([]T, len(v.data))

	for i, val := range v.data {
		newData[i] = cb(i, val)
	}

	return NewVec[T](newData...)
}

// Filter iterates over the vector and calls the callback function for each value.
// The callback function should return true if the value should be kept, false otherwise.
// It returns a new vector with the values that returned true.
func (v *Vec[T]) Filter(cb func(int, T) bool) *Vec[T] {
	newData := make([]T, 0)

	for i, val := range v.data {
		if cb(i, val) {
			newData = append(newData, val)
		}
	}

	return NewVec[T](newData...)
}

// Reduce iterates over the vector and calls the callback function for each value.
// The callback function should return the new value for the accumulator.
// It returns the final value of the accumulator.
func (v *Vec[T]) Reduce(cb func(T, T) T) T {
	var acc T

	for _, curr := range v.data {
		acc = cb(acc, curr)
	}

	return acc
}

// Reverse reverses the vector.
func (v *Vec[T]) Reverse() *Vec[T] {
	for i := 0; i < len(v.data)/2; i++ {
		v.data[i], v.data[len(v.data)-i-1] = v.data[len(v.data)-i-1], v.data[i]
	}

	return v
}

// ToSlice returns the vector as a slice.
func (v *Vec[T]) ToSlice() []T {
	return v.data
}
