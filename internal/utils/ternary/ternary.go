package ternary

import "strings"

// ValueOrDefault returns the input value if the conditionFunc is true; otherwise, it returns the specified falseValue.
func ValueOrDefault[T any](value T, conditionFunc func(T) bool, falseValue T) T {
	if conditionFunc(value) {
		return value
	}
	return falseValue
}

// NotNilOrWhitespace checks if the given string is not empty and does not consist solely of whitespace characters.
func NotNilOrWhitespace(value string) bool {
	return value != "" && strings.TrimSpace(value) != ""
}
