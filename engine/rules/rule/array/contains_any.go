package rule_array

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type ArrayContainsAny struct {
	Type  string        `json:"type"`
	Field string        `json:"field"`
	Value []interface{} `json:"value"`
	hash  *string
}

func (c *ArrayContainsAny) GetType() string {
	return "rules.array.contains_any"
}

func NewArrayContainsAny() rules.RulesItf {
	o := new(ArrayContainsAny)
	o.Type = o.GetType()
	return o
}

func (c *ArrayContainsAny) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.Type, c.Field, c.Value)
		c.hash = &hash
	}
	return *c.hash
}

func (c *ArrayContainsAny) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
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

	for _, val := range *intv {
		for _, check := range c.Value {
			if val == check {
				return true, nil
			}
		}
	}

	return false, nil
}
