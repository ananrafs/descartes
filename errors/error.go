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

	ErrNilEvaluator = func(slug string) error {
		return fmt.Errorf("evaluator is nil for %s", slug)
	}

	ErrInvalidOnCreation = func(slug string) error {
		return fmt.Errorf("failed to create %s", slug)
	}
)
