package rule_array

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
	json "github.com/json-iterator/go"
	"reflect"
)

type ArrayContainsStruct struct {
	Type  string      `json:"type"`
	Field string      `json:"field"`
	Value interface{} `json:"value"`
	mVal  map[string]interface{}
	hash  *string
}

func (c *ArrayContainsStruct) GetType() string {
	return "rules.array.contains.struct"
}

func NewArrayContainsStruct() rules.RulesItf {
	o := new(ArrayContainsStruct)
	o.Type = o.GetType()
	o.mVal = make(map[string]interface{})
	return o
}

func (c *ArrayContainsStruct) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.Type, c.Field, c.Value)
		c.hash = &hash
	}
	return *c.hash
}

func (c *ArrayContainsStruct) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
	if ok := facts.GetCacheInstance().TryGet(c.GetHash(), &isMatch); ok {
		return isMatch, nil
	}
	defer func() {
		facts.GetCacheInstance().SetCache(c.GetHash(), isMatch)
	}()
	param := facts.GetMap()
	v, ok := param[c.Field]
	if !ok {
		return false, common.ErrorNotFoundOnMap(c.Field)
	}

	intv := new([]interface{})
	if err = common.ConvertToArray(v, intv); err != nil {
		return false, err
	}

parentLoop:
	for _, val := range *intv {
		factField, ok := val.(map[string]interface{})
		if !ok {
			return false, common.ErrorCasting(val)
		}

		for k, v := range c.mVal {
			// read attribute on field-k on fact
			factFieldAttribute, ok := factField[k]
			if !ok {
				//continue to next element on fact array
				// since there are no attribute define in rule on this array struct
				continue parentLoop
			}

			ruleMap, ok := v.(map[interface{}]bool)
			if ok {
				// read value on rule from transformed array to map
				_, ok := ruleMap[factFieldAttribute]
				if !ok {
					//continue to next element on fact array
					// since there are no value equal on rule
					continue parentLoop
				}
			} else {
				if factFieldAttribute != v {
					//continue to next element on fact array
					// since this rule value isnt array and not equal with given array struct
					continue parentLoop
				}
			}
		}
		return true, nil
	}

	return false, nil
}

func (o *ArrayContainsStruct) UnmarshalJSON(data []byte) (err error) {
	var m map[string]json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}

	for k, val := range m {
		switch k {
		case "type":
			if err := json.Unmarshal(val, &o.Type); err != nil {
				return err
			}
		case "field":
			if err := json.Unmarshal(val, &o.Field); err != nil {
				return err
			}
		case "value":
			if err := json.Unmarshal(val, &o.Value); err != nil {
				return err
			}
			cVal, ok := o.Value.(map[string]interface{})
			if !ok {
				return common.ErrorCasting(o.Value)
			}
			for k, v := range cVal {
				val := v
				if refVal := reflect.ValueOf(v); refVal.Kind() == reflect.Slice {
					innerMap := map[interface{}]bool{}
					for i := 0; i < refVal.Len(); i++ {
						innerMap[refVal.Index(i).Interface()] = true
					}
					val = innerMap
				}
				o.mVal[k] = val
			}
		}
	}

	return
}
