package common

import (
	"reflect"
	"strconv"
	"strings"
)

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}
type NumberConverter[T Number] func(source interface{}, dest *T) error

func (nc NumberConverter[T]) WithFromMap(mp map[string]interface{}) NumberConverter[T] {
	return func(source interface{}, dest *T) error {
		numField := new(string)
		// check if its using template
		if match := ParseFromTemplate(source, numField); match {
			val, err := LookUpMap(mp, *numField)
			if err == nil {
				return nc(val, dest)
			}
		}

		return nc(source, dest)
	}
}

func Convert[T Number]() (nc NumberConverter[T]) {
	return ConvertNumber[T]
}

func ConvertNumber[T Number](source interface{}, dest *T) error {
	_valueRef := reflect.ValueOf(source)
	switch _valueRef.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		*dest = T(_valueRef.Int())
	case reflect.Float32, reflect.Float64:
		*dest = T(_valueRef.Float())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		*dest = T(_valueRef.Uint())
	case reflect.String:
		return StringConverter[T](_valueRef.String(), dest)
	default:
		return ErrorCasting(source)
	}
	return nil
}

func StringConverter[T Number](source string, dest *T) error {
	temp, err := strconv.Atoi(source)
	if err == nil {
		*dest = T(temp)
		return nil
	}
	tempF, err := strconv.ParseFloat(source, 64)
	if err != nil {
		return ErrorCasting(source)
	}
	*dest = T(tempF)
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
