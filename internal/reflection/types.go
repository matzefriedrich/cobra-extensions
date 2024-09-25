package reflection

import (
	"reflect"
)

type valueItem struct {
	value     reflect.Value
	valueType reflect.Type
}
