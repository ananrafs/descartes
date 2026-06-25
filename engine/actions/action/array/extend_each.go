package array

import (
	"fmt"
	"reflect"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/actions"
	"github.com/ananrafs/descartes/engine/facts"
)

type ExtendEach struct {
	Type   string                 `json:"type"`
	Field  string                 `json:"field"`
	Object map[string]interface{} `json:"object"`
}

func (a *ExtendEach) GetType() string {
	return "actions.array.extend.each"
}

func NewExtendEach() actions.ActionsItf {
	o := new(ExtendEach)
	o.Type = o.GetType()
	o.Object = make(map[string]interface{})
	return o
}

func (a *ExtendEach) Do(facts facts.FactsItf) (res interface{}, err error) {
	param := facts.GetMap()
	field := a.Field
	var fieldItf interface{}
	if match := common.DeepTemplateEvaluateFromMap(param, field, &fieldItf); match {
		field = fmt.Sprintf("%v", fieldItf)
	}

	srcCopy := common.ManipulateMap(map[string]interface{}{}).DeepCopy(a.Object)

	objValue, ok := param[field]
	if !ok {
		objValue = map[string]interface{}{}
	}

	refValue := reflect.ValueOf(objValue)
	switch refValue.Kind() {
	case reflect.Slice:
		for i := 0; i < refValue.Len(); i++ {
			item := refValue.Index(i)
			if !item.IsValid() {
				return nil, common.ErrorCasting(item)
			}
			err := appendIteratedObject(item, param, srcCopy)
			if err != nil {
				return nil, err
			}
		}
	case reflect.Map:
		for _, key := range refValue.MapKeys() {
			item := refValue.MapIndex(key)
			if !item.IsValid() {
				return nil, common.ErrorCasting(item)
			}
			err := appendIteratedObject(item, param, srcCopy)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, common.ErrorCasting(objValue)
	}

	param[field] = objValue
	return
}

func appendIteratedObject(item reflect.Value, param, srcCopy map[string]interface{}) error {
	iteratedItem, ok := item.Interface().(map[string]interface{})
	if !ok {
		return common.ErrorCasting(item)
	}

	appendObject(param, srcCopy, iteratedItem)
	return nil
}

func appendObject(param, src, dest map[string]interface{}) {
	for key, value := range src {
		var keyMapField interface{}
		if match := common.DeepTemplateEvaluateFromMap(param, key, &keyMapField); match {
			key = fmt.Sprintf("%v", keyMapField)
		}

		_map := map[string]interface{}{}
		isObj := common.ExtractMap(value, &_map, func(_key *string, _val *interface{}) {
			var valMapField interface{}
			if match := common.DeepTemplateEvaluateFromMap(param, *_val, &valMapField); match {
				*_val = valMapField
			}
		})
		if isObj {
			if currentMapValue, ok := dest[key]; ok {
				existingMap, ok := currentMapValue.(map[string]interface{})
				if ok {
					dest[key] = common.ManipulateMap(_map).Merge(existingMap)
					continue
				}
			}
			dest[key] = _map
			continue
		}

		var valMapField = value
		common.DeepTemplateEvaluateFromMap(param, value, &valMapField)
		dest[key] = valMapField
	}
}
