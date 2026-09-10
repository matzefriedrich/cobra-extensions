package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Stack_IsEmpty_returns_true_for_an_empty_stack(t *testing.T) {
	// Arrange
	sut := MakeStack[int]()

	// Act
	actual := sut.IsEmpty()

	// Assert
	assert.True(t, actual)
}

func Test_Stack_IsEmpty_returns_false_after_values_are_pushed(t *testing.T) {
	// Arrange
	sut := MakeStack[string]()
	sut.Push("value")

	// Act
	actual := sut.IsEmpty()

	// Assert
	assert.False(t, actual)
}

func Test_Stack_Any_returns_false_for_an_empty_stack(t *testing.T) {
	// Arrange
	sut := MakeStack[int]()

	// Act
	actual := sut.Any()

	// Assert
	assert.False(t, actual)
}

func Test_Stack_Any_returns_true_for_a_non_empty_stack(t *testing.T) {
	// Arrange
	sut := MakeStack[int]()
	sut.Push(42)

	// Act
	actual := sut.Any()

	// Assert
	assert.True(t, actual)
}

func Test_Stack_Push_appends_values_to_the_stack(t *testing.T) {
	// Arrange
	sut := MakeStack[int]()

	// Act
	sut.Push(1)
	sut.Push(2)

	// Assert
	assert.Equal(t, []int{1, 2}, []int(sut))
}

func Test_Stack_Push_is_variadic_and_appends_all_values(t *testing.T) {
	// Arrange
	sut := MakeStack[int]()

	// Act
	sut.Push(1, 2, 3)

	// Assert
	assert.Equal(t, []int{1, 2, 3}, []int(sut))
}

func Test_Stack_Pop_returns_the_top_item_and_removes_it(t *testing.T) {
	// Arrange
	sut := MakeStack[int]()
	sut.Push(1)
	sut.Push(2)

	// Act
	actual := sut.Pop()

	// Assert
	assert.Equal(t, 2, actual)
	assert.Equal(t, []int{1}, []int(sut))
}

func Test_Stack_Pop_returns_items_in_lifo_order(t *testing.T) {
	// Arrange
	sut := MakeStack[string]()
	sut.Push("first")
	sut.Push("second")
	sut.Push("third")

	// Act
	third := sut.Pop()
	second := sut.Pop()
	first := sut.Pop()

	// Assert
	assert.Equal(t, "third", third)
	assert.Equal(t, "second", second)
	assert.Equal(t, "first", first)
	assert.True(t, sut.IsEmpty())
}
