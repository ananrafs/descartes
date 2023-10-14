package actions

import (
	"github.com/ananrafs/descartes/common"
)

var (
	actionsMap map[string]ActionsItf = make(map[string]ActionsItf)
)

type ActionsItf interface {
	common.TypeCheckerItf
	New() ActionsItf
	Do(map[string]interface{}) (interface{}, error)
}

func Init(acts ...ActionsItf) {
	for _, action := range acts {
		actionsMap[action.GetType()] = action
	}
}

func Get(rulesType string) ActionsItf {
	act, ok := actionsMap[rulesType]
	if ok {
		return act
	}

	return nil
}
