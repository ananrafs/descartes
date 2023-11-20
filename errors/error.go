package errors

import (
	"fmt"
)

var (
	ErrLawNotFound = func(slug string) error {
		return fmt.Errorf("law %s not found", slug)
	}

	ErrFactsNotMatch = func(slug string) error {
		return fmt.Errorf("not match on %s law", slug)
	}
)
