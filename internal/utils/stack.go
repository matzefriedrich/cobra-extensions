package utils

// Stack represents a generic stack data structure.
type Stack[TValue any] []TValue

// MakeStack creates and returns an empty Stack of the specified generic type TValue.
func MakeStack[TValue any]() Stack[TValue] {
	return Stack[TValue]{}
}

// Any returns true if the stack has items, otherwise false.
func (s *Stack[TValue]) Any() bool {
	return !s.IsEmpty()
}

// IsEmpty returns true if the stack is currently empty, otherwise false.
func (s *Stack[TValue]) IsEmpty() bool {
	return len(*s) == 0
}

// Push adds one or more values of type TValue to the top of the stack.
func (s *Stack[TValue]) Push(value ...TValue) {
	*s = append(*s, value...)
}

// Pop removes and returns the top item from the stack.
func (s *Stack[TValue]) Pop() TValue {
	lastItemIndex := len(*s) - 1
	element := (*s)[lastItemIndex]
	*s = (*s)[:lastItemIndex]
	return element
}
