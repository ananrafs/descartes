package common

import (
	"reflect"
	"regexp"
	"strings"
)

var (
	templateRegex = regexp.MustCompile(`\{\{([^{}]+)\}\}`)
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
	intf, ok := source.(bool)
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

	target := templateRegex.FindStringSubmatch(strSource)
	if len(target) < 1 {
		return false
	}
	*dest = strings.TrimSpace(target[1])
	return true
}

func CopyMap(src map[string]interface{}) (dest map[string]interface{}) {
	dest = make(map[string]interface{})
	for k, v := range src {
		dest[k] = v
	}
	return
}

func LookupMap(mp map[string]interface{}, index int, source []string) (interface{}, error) {
	if val, ok := mp[source[index]]; ok {
		childMap, ok := val.(map[string]interface{})
		if !ok {
			return val, nil
		}
		index++
		return LookupMap(childMap, index, source)
	}
	return nil, ErrorNotFoundOnMap(source[index])
}
