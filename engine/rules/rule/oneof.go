package rule

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
	json "github.com/json-iterator/go"
)

type OneOf struct {
	Type  string      `json:"type"`
	Field string      `json:"field"`
	Value interface{} `json:"value"`

	mValue map[interface{}]bool
	hash   *string
}

func (o *OneOf) GetType() string {
	return "rules.oneof"
}

func NewOneOf() rules.RulesItf {
	o := new(OneOf)
	o.Type = o.GetType()
	o.mValue = make(map[interface{}]bool)
	return o
}

func (o *OneOf) GetHash() string {
	for o.hash == nil {
		hash := common.CreateHash(o.Type, o.Field, o.Value)
		o.hash = &hash
	}
	return *o.hash
}

func (o *OneOf) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
	if ok := facts.GetCacheInstance().TryGet(o.GetHash(), &isMatch); ok {
		return isMatch, nil
	}
	defer func() {
		facts.GetCacheInstance().SetCache(o.GetHash(), isMatch)
	}()
	param := facts.GetMap()
	v, ok := param[o.Field]
	if !ok {
		return false, common.ErrorNotFoundOnMap(o.Field)
	}

	return o.mValue[v], nil
}

func (o *OneOf) UnmarshalJSON(data []byte) (err error) {
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
			intv := new([]interface{})
			if err = common.ConvertToArray(o.Value, intv); err != nil {
				continue
			}
			for _, val := range *intv {
				o.mValue[val] = true
			}
		}
	}

	return
}
