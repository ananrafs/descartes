package core

import (
	"github.com/ananrafs/descartes/errors"
	"github.com/ananrafs/descartes/law"
)

var (
	lawDictionary = make(map[string]law.Law)
)

func Eval(fact law.Fact) (interface{}, error) {
	lawSelected, ok := lawDictionary[fact.Slug]
	if !ok {
		return nil, errors.ErrLawNotFound(fact.Slug)
	}

	return lawSelected.Judge(fact.Facts)
}
