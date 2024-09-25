package utils

// Stack represents a generic stack data structure.
type Stack[TValue any] []TValue

// MakeStack creates and returns an empty Stack of the specified generic type TValue.
func MakeStack[TValue any]() Stack[TValue] {
	return Stack[TValue]{}
}

func (s *Stack[TValue]) emptyValue() TValue {
	var value TValue
	return value
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

// TryPop attempts to remove and return the top item from the stack.
// If the stack is empty, it returns the default value and false.
func (s *Stack[TValue]) TryPop() (TValue, bool) {
	if s.IsEmpty() {
		return s.emptyValue(), false
	} else {
		return s.Pop(), true
	}
}

// Pop removes and returns the top item from the stack.
func (s *Stack[TValue]) Pop() TValue {
	lastItemIndex := len(*s) - 1
	element := (*s)[lastItemIndex]
	*s = (*s)[:lastItemIndex]
	return element
}

// Peek returns the top item from the stack without removing it.
func (s *Stack[TValue]) Peek() TValue {
	lastItemIndex := len(*s) - 1
	element := (*s)[lastItemIndex]
	return element
}

// TryPeek attempts to return the top item from the stack without removing it.
// Returns the item and true if the stack is not empty, otherwise it returns the default value and false.
func (s *Stack[TValue]) TryPeek() (TValue, bool) {
	if s.IsEmpty() {
		return s.emptyValue(), false
	} else {
		return s.Peek(), true
	}
}
