package libs

import (
	"fmt"
	"reflect"
	"strings"
)

func StructToMapCmp(in interface{}, tag string, changemap map[string]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		if tag == "" {
			out[fi.Name] = v.Field(i).Interface()
		} else {
			if tagv := fi.Tag.Get(tag); tagv != "" {
				// set key of map to value in struct field
				if tag == "edit" && changemap != nil {
					_, have := changemap[tagv]
					if have {
						out[tagv] = v.Field(i).Interface()
					}
				} else {
					out[tagv] = v.Field(i).Interface()
				}

			}
		}

	}
	return out, nil
}

func StructToMap(in interface{}, tag string) (map[string]interface{}, error) {
	return StructToMapCmp(in, tag, nil)
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
