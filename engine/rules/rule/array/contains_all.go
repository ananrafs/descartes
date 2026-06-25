package rule_array

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type ArrayContainsAll struct {
	Type  string        `json:"type"`
	Field string        `json:"field"`
	Value []interface{} `json:"value"`
	hash  *string
}

func (c *ArrayContainsAll) GetType() string {
	return "rules.array.contains_all"
}

func NewArrayContainsAll() rules.RulesItf {
	o := new(ArrayContainsAll)
	o.Type = o.GetType()
	return o
}

func (c *ArrayContainsAll) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.Type, c.Field, c.Value)
		c.hash = &hash
	}
	return *c.hash
}

func (c *ArrayContainsAll) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
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

	// Create a map for O(1) lookups
	factMap := make(map[interface{}]bool)
	for _, val := range *intv {
		factMap[val] = true
	}

	// Check if all values in c.Value are present in factMap
	for _, check := range c.Value {
		if !factMap[check] {
			return false, nil
		}
	}

	return true, nil
}
