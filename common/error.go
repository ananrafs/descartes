package common

import (
	"fmt"
)

var ()

func ErrorCasting(source interface{}) error {
	return fmt.Errorf("failed to cast %s", source)
}

func ErrorNotFoundOnMap(field string) error {
	return fmt.Errorf("not found on field for %s", field)

}

func ErrorOutOfBound(field string) error {
	return fmt.Errorf("index out of bound for %s", field)
}
