package core

import (
	"fmt"

	"github.com/ananrafs/descartes/law"
)

var (
	lawDictionary = make(map[string]law.Law)
)

func Eval(fact law.Fact) (interface{}, error) {
	lawSelected, ok := lawDictionary[fact.Slug]
	if !ok {
		return nil, fmt.Errorf("law %s not found", fact.Slug)
	}

	return lawSelected.Judge(fact.Param)

}
