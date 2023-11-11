package action_map

import (
	"fmt"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/actions"
	"github.com/ananrafs/descartes/engine/facts"
)

type Append struct {
	Type   string                 `json:"type"`
	Field  string                 `json:"field"`
	Object map[string]interface{} `json:"object"`
}

func (a *Append) GetType() string {
	return "actions.map.append"
}

func (a *Append) New() actions.ActionsItf {
	o := new(Append)
	o.Type = o.GetType()
	o.Object = make(map[string]interface{})
	return o
}

func (a Append) Do(facts facts.FactsItf) (res interface{}, err error) {
	param := facts.GetMap()
	field := a.Field
	var fieldItf interface{}
	if match := common.DeepTemplateEvaluateFromMap(param, field, &fieldItf); match {
		field = fmt.Sprintf("%v", fieldItf)
	}

	objValue, ok := param[field]
	if !ok {
		objValue = map[string]interface{}{}
	}

	objMapped, ok := objValue.(map[string]interface{})
	if !ok {
		return
	}

	// ObjectToAppend := common.CopyMap(a.Object)

	for key, value := range a.Object {
		var keyMapField interface{}
		if match := common.DeepTemplateEvaluateFromMap(param, key, &keyMapField); match {
			key = fmt.Sprintf("%v", keyMapField)
		}

		_map := map[string]interface{}{}
		isObj := common.RecurringMap(value, &_map, func(_key *string, _val *interface{}) {
			var valMapField interface{}
			if match := common.DeepTemplateEvaluateFromMap(param, *_val, &valMapField); match {
				*_val = valMapField
			}
		})
		if isObj {
			if currentMapValue, ok := objMapped[key]; ok {
				existingMap, ok := currentMapValue.(map[string]interface{})
				if ok {
					objMapped[key] = common.MergeMap(existingMap, _map)
					continue
				}
			}
			objMapped[key] = _map
			continue
		}

		var valMapField = value
		common.DeepTemplateEvaluateFromMap(param, value, &valMapField)
		objMapped[key] = valMapField
	}
	param[field] = objMapped
	return
}
