package action

import (
	"fmt"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/actions"
	"github.com/ananrafs/descartes/engine/facts"
)

type Action map[string]interface{}

func (c Action) GetType() string {
	return ""
}

func NewAction() actions.ActionsItf {
	newAction := make(Action)
	return &newAction
}

func (c Action) Do(facts facts.FactsItf) (res interface{}, err error) {
	param := facts.GetMap()
	response := common.CopyMap(map[string]interface{}(c))

	for key, value := range c {

		var keyMapField interface{}
		if match := common.DeepTemplateEvaluateFromMap(param, key, &keyMapField); match {
			defer func(key string) {
				delete(response, key)
			}(key)

			key = fmt.Sprintf("%v", keyMapField)
		}

		var valMapField interface{}
		if match := common.DeepTemplateEvaluateFromMap(param, value, &valMapField); match {
			value = valMapField
		}
		response[key] = value

	}

	return response, nil
}
