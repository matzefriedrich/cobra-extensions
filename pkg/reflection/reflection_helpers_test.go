package reflection

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type TestStruct struct {
	Name     string
	Age      int64
	Remember bool
}

func TestReflectedMember_EnumerateFields_fields_contain_expected_values(t *testing.T) {

	// Arrange
	n := &TestStruct{
		Age:      99,
		Remember: true,
		Name:     "John Doe",
	}

	sut := ReflectObject(n)

	const expectedNumFields = 3
	numFields := 0

	actualName := ""
	actualAge := int64(0)
	actualRemember := false

	// Act
	sut.EnumerateFields(func(index int, field ReflectedField) {
		numFields++
		switch field.typeKind() {
		case reflect.String:
			actualName = field.value.String()
		case reflect.Int64:
			actualAge = field.value.Int()
		case reflect.Bool:
			actualRemember = field.value.Bool()
		default:
			panic(fmt.Sprintf("unexpected field type: %v", field.typeKind()))
		}
	})

	// Assert
	assert.Equal(t, expectedNumFields, numFields)
	assert.Equal(t, n.Name, actualName)
	assert.Equal(t, n.Age, actualAge)
	assert.Equal(t, n.Remember, actualRemember)
}
