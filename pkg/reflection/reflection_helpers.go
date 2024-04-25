package reflection

import "reflect"

type ReflectedObject struct {
	instanceValue reflect.Value
	objectType    reflect.Type
}

type ReflectedField struct {
	field reflect.StructField
	value reflect.Value
}

func ReflectObject(n interface{}) *ReflectedObject {
	value := reflect.ValueOf(n)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	valueType := reflect.TypeOf(value.Interface())
	return &ReflectedObject{instanceValue: value, objectType: valueType}
}

func (m *ReflectedObject) Kind() reflect.Kind {
	return m.objectType.Kind()
}

type FieldEnumeratorCallback func(index int, field ReflectedField)

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
