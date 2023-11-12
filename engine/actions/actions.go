package actions

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
)

var (
	actionsMap map[string]Factory = make(map[string]Factory)
)

type Factory func() ActionsItf

type ActionsItf interface {
	common.TypeCheckerItf
	Do(facts.FactsItf) (interface{}, error)
}

func Init(acts ...Factory) {
	for _, action := range acts {
		actionsMap[action().GetType()] = action
	}
}

func Get(rulesType string) ActionsItf {
	act, ok := actionsMap[rulesType]
	if ok {
		return act()
	}

	return nil
}

func GetCatalog() []Factory {
	res := make([]Factory, 0, len(actionsMap))
	for _, act := range actionsMap {
		res = append(res, act)
	}

	return res
}
