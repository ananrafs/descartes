package common

import (
	"regexp"
	"strings"
)

var (
	handlebars = regexp.MustCompile(`\{\{([^{}]+)\}\}`)
)

func ConvertToInt(source interface{}, dest *int) error {
	intf, ok := source.(float64)
	if !ok {
		return ErrorCasting(source)
	}

	*dest = int(intf)
	return nil
}

func ParseFromHandlebars(source interface{}, dest *string) (isMatch bool) {
	strSource, ok := source.(string)
	if !ok {
		return false
	}

	target := handlebars.FindStringSubmatch(strSource)
	if len(target) < 1 {
		return false
	}
	*dest = strings.TrimSpace(target[1])
	return true
}
