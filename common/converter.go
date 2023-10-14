package common

import (
	"regexp"
	"strings"
)

var (
	templateRegex = regexp.MustCompile(`\{\{([^{}]+)\}\}`)
)

func ConvertToInt(source interface{}, dest *int) error {
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
