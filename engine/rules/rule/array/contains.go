package rule_array

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type ArrayContains struct {
	RuleType string      `json:"type"`
	Field    string      `json:"field"`
	Value    interface{} `json:"value"`
	hash     *string
}

func (c *ArrayContains) GetType() string {
	return "rules.array.contains"
}

func (c *ArrayContains) New() rules.RulesItf {
	return new(ArrayContains)
}

func (c *ArrayContains) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType, c.Field, c.Value)
		c.hash = &hash
	}
	return *c.hash
}

func (c *ArrayContains) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
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
		if val == c.Value {
			return true, nil
		}
	}

	return false, nil
}
