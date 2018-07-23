package template

import (
	"fmt"
	"reflect"
	"strings"
)

func join(in []interface{}, sepprefixsuffix ...string) string {

	_len := len(sepprefixsuffix)

	if _len != 1 && _len != 3 {
		panic("Function requires one spearator or one separator, prefix, and suffix")
	}

	items := make([]string, len(in))

	for i, v := range in {
		if _len == 3 {
			items[i] = fmt.Sprintf("%s%s%s", sepprefixsuffix[1], v, sepprefixsuffix[2])
		} else {
			items[i] = fmt.Sprintf("%s", v)
		}
	}

	return strings.Join(items, sepprefixsuffix[0])
}

func project(in []interface{}, property string) []interface{} {

	var result []interface{}

	var appendValue func(list []interface{}, v reflect.Value)

	appendValue = func(list []interface{}, v reflect.Value) {
		kind := v.Kind()
		if kind == reflect.Ptr {
			appendValue(list, v.Elem())
		} else if kind == reflect.Struct {
			result = append(result, v.FieldByName(property).Interface())
		} else if kind == reflect.Map {
			result = append(result, v.Interface().(map[string]interface{})[property])
		} else {
			panic("Value must be a struct or map")
		}
	}

	for _, el := range in {
		v := reflect.ValueOf(el)
		appendValue(result, v)
	}

	return result

}
