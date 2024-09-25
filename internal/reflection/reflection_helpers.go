package reflection

import "reflect"

// ReflectedObject is a struct that wraps reflection information for a given object instance.
type ReflectedObject struct {
	instanceValue reflect.Value
	objectType    reflect.Type
}

// ReflectedField represents a field within a struct that includes both its reflective type and value information.
type ReflectedField struct {
	field reflect.StructField
	value reflect.Value
}

// ReflectObject takes an interface and returns a reference to a ReflectedObject, containing reflection info about the instance.
func ReflectObject(n interface{}) *ReflectedObject {
	value := reflect.ValueOf(n)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	valueType := reflect.TypeOf(value.Interface())
	return &ReflectedObject{instanceValue: value, objectType: valueType}
}

// Kind returns the reflection kind of the wrapped object type.
func (m *ReflectedObject) Kind() reflect.Kind {
	return m.objectType.Kind()
}

// FieldEnumeratorCallback is a type for a callback function used to iterate over fields of a struct.
type FieldEnumeratorCallback func(index int, field ReflectedField)

// EnumerateFields iterates over each field of the underlying struct, invoking the provided callback for each field.
func (m *ReflectedObject) EnumerateFields(iterFunc FieldEnumeratorCallback) {
	if m.Kind() != reflect.Struct {
		return
	}
	numFields := m.objectType.NumField()
	for i := 0; i < numFields; i++ {
		structField := m.objectType.Field(i)
		structFieldValue := m.instanceValue.Field(i)
		field := ReflectedField{field: structField, value: structFieldValue}
		iterFunc(i, field)
	}
}

func (f *ReflectedField) typeKind() reflect.Kind {
	return f.field.Type.Kind()
}

func (f *ReflectedField) isType(t any) bool {
	return f.field.Type == reflect.TypeOf(t)
}

func (f *ReflectedField) getInterfaceValue() interface{} {
	return f.value.Interface()
}
