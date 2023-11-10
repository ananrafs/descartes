package common

import (
	"reflect"
	"strconv"
	"strings"
)

func ConvertToInt(source interface{}, dest *int) error {
	_valueRef := reflect.ValueOf(source)
	switch _valueRef.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		*dest = int(_valueRef.Int())
	case reflect.Float32, reflect.Float64:
		*dest = int(_valueRef.Float())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		*dest = int(_valueRef.Uint())
	case reflect.String:
		temp, err := strconv.Atoi(_valueRef.String())
		if err != nil {
			return ErrorCasting(source)
		}
		*dest = int(temp)
	default:
		return ErrorCasting(source)
	}

	return nil
}

type IntConverter func(interface{}, *int) error

func ConvertInt() IntConverter {
	return ConvertToInt
}

func (ic IntConverter) WithFromMap(mp map[string]interface{}) IntConverter {
	return func(source interface{}, dest *int) error {
		numField := new(string)
		// check if its using template
		if match := ParseFromMustacheTemplate(source, numField); match {
			var ok bool
			source, ok = mp[*numField]
			if !ok {
				return ic(source, dest)
			}
		}

		return ic(source, dest)
	}
}

type FloatConverter func(interface{}, *float64) error

func ConvertToFloat(source interface{}, dest *float64) error {
	_valueRef := reflect.ValueOf(source)
	switch _valueRef.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		*dest = float64(_valueRef.Int())
	case reflect.Float32, reflect.Float64:
		*dest = float64(_valueRef.Float())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		*dest = float64(_valueRef.Uint())
	case reflect.String:
		temp, err := strconv.ParseFloat(_valueRef.String(), 64)
		if err != nil {
			return ErrorCasting(source)
		}
		*dest = temp
	default:
		return ErrorCasting(source)
	}

	return nil
}

func ConvertFloat() FloatConverter {
	return ConvertToFloat
}

func (ic FloatConverter) WithFromMap(mp map[string]interface{}) FloatConverter {
	return func(source interface{}, dest *float64) error {
		numField := new(string)
		// check if its using template
		if match := ParseFromMustacheTemplate(source, numField); match {
			dotSplittedSlice := strings.Split(*numField, ".")
			val, err := LookupMap(mp, 0, dotSplittedSlice)
			if err == nil {
				return ic(val, dest)
			}

		}

		return ic(source, dest)
	}
}

func ConvertFromMap(source interface{}, dest *int) error {
	intf, ok := source.(float64)
	if !ok {
		return ErrorCasting(source)
	}

	*dest = int(intf)
	return nil
}

func ConvertToBool(source interface{}, dest *bool) error {
	if reflect.ValueOf(source).Kind() == reflect.String {
		sourceTemp := strings.ToLower(reflect.ValueOf(source).String())
		temp := sourceTemp == "true"
		if !temp && source != "false" {
			return ErrorCasting(source)
		}
		*dest = temp
		return nil
	}
	intf, ok := source.(bool)
	if !ok {
		return ErrorCasting(source)
	}

	*dest = intf
	return nil
}

func ConvertToArray(source interface{}, dest *[]interface{}) error {
	intf, ok := source.([]interface{})
	if !ok {
		return ErrorCasting(source)
	}

	*dest = intf
	return nil
}

func ParseFromMustacheTemplate(source interface{}, dest *string) (isMatch bool) {
	strSource, ok := source.(string)
	if !ok {
		return false
	}

	return GetTemplatedString(strSource, dest)
}

func GetTemplatedString(source string, dest *string) bool {
	open := "{{"
	close := "}}"

	if source[:2] == open && source[len(source)-2:] == close {
		trimmed := strings.TrimSpace(source[2 : len(source)-2])
		*dest = trimmed
		return true
	}
	return false

}
func DeepParseFromMustacheTemplate(source interface{}, dest *string) (isMatch bool, depth int) {
	strSource, ok := source.(string)
	if !ok {
		return false, 0
	}

	depth = 0
	// Find all matches in the input string
	for {
		match := GetTemplatedString(strSource, dest)
		if !match {
			break
		}

		strSource = *dest
		depth++
	}
	if depth == 0 {
		return false, depth
	}

	*dest = strSource

	return true, depth
}

func DeepTemplateEvaluateFromMap(mp map[string]interface{}, src interface{}, dest *interface{}) bool {
	keyMapField := ""

	// check if its using template
	if match, deep := DeepParseFromMustacheTemplate(src, &keyMapField); match {
		var (
			valueField interface{}
			lookUpMap  map[string]interface{} = mp
			ok         bool
		)

		for i := 0; i < deep; i++ {
			if valueField != nil {
				keyMapField, ok = valueField.(string)
				if !ok {
					break
				}
			}

			valueField, _ = LookUpRecursiveMap(lookUpMap, keyMapField)
			if valueField == nil {

				return false
			}

		}

		*dest = valueField

		return true
	}

	return false
}

func CopyMap(src map[string]interface{}) (dest map[string]interface{}) {
	dest = make(map[string]interface{})
	for k, v := range src {
		dest[k] = v
	}
	return
}

func LookUpRecursiveMap(mp map[string]interface{}, src string) (interface{}, error) {
	return LookupMap(mp, 0, strings.Split(src, "."))
}

func LookupMap(mp map[string]interface{}, index int, source []string) (res interface{}, err error) {
	if index >= len(source) {
		return nil, nil
	}
	if val, ok := mp[source[index]]; ok {
		childMap, ok := val.(map[string]interface{})
		if !ok {
			return val, nil
		}
		index++
		lookup, _ := LookupMap(childMap, index, source)
		if lookup != nil {
			val = lookup
		}
		return val, nil
	}
	return nil, ErrorNotFoundOnMap(source[index])
}

// RecurringMap is method to evaluate interface{}.
// if source is map, then it will lookup to key-pairs and perform recursive.
// modifiers is used to alter or modify value of map
func RecurringMap(source interface{}, dest *map[string]interface{}, modifiers ...func(*string, *interface{})) bool {
	sMap, ok := source.(map[string]interface{})
	if !ok {
		return false
	}

	copiedMap := CopyMap(sMap)
	for key, val := range copiedMap {
		var _dest map[string]interface{}
		isObj := RecurringMap(val, &_dest, modifiers...)
		if isObj {
			copiedMap[key] = _dest
			continue
		}
		for _, mod := range modifiers {
			mod(&key, &val)
		}
		copiedMap[key] = val
	}
	*dest = copiedMap

	return true
}
