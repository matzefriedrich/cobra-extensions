package reflection

import "reflect"

type ReflectedMember struct {
	value      reflect.Value
	memberType reflect.Type
}

type ReflectedField struct {
	field reflect.StructField
	value reflect.Value
}

func (m *ReflectedMember) Kind() reflect.Kind {
	return m.memberType.Kind()
}

type FieldEnumeratorCallback func(index int, field ReflectedField)

func (m *ReflectedMember) EnumerateFields(iterFunc FieldEnumeratorCallback) {
	if m.Kind() != reflect.Struct {
		return
	}
	numFields := m.memberType.NumField()
	for i := 0; i < numFields; i++ {
		structField := m.memberType.Field(i)
		structFieldValue := m.value.Field(i)
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
