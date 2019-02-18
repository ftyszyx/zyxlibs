package libs

import (
	"reflect"
	"strings"
)

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[strings.ToLower(t.Field(i).Name)] = v.Field(i).Interface()
	}
	return data
}

//按照结构体清除map无用字段,只有在结构体中有的字段，才会放进去
func ClearMapByStructTag(data map[string]interface{}, in interface{}, tag string) map[string]interface{} {
	out := make(map[string]interface{})
	v := reflect.ValueOf(in)
	typ := v.Type()
	if data == nil {
		return data
	}
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		lowName := strings.ToLower(fi.Name)
		if tag == "" {
			_, have := data[lowName]
			if have {
				out[lowName] = data[lowName]
			}
		} else {
			if tagv := fi.Tag.Get(tag); tagv != "" {
				_, have := data[tagv]
				if have {
					out[tagv] = data[tagv]
				}

			}
		}

	}
	return out

}

func ClearMapByStruct(data map[string]interface{}, in interface{}) map[string]interface{} {
	return ClearMapByStructTag(data, in, "")
}
