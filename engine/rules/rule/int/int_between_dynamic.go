package rule_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleIntBetweenDynamic struct {
	RuleType string      `json:"type"`
	Mid      interface{} `json:"mid"`
	Start    interface{} `json:"start"`
	End      interface{} `json:"end"`
	hash     *string
}

func (c *RuleIntBetweenDynamic) GetType() string {
	return "rules.int.between.dynamic"
}

func (c *RuleIntBetweenDynamic) New() rules.RulesItf {
	return new(RuleIntBetweenDynamic)
}

func (c *RuleIntBetweenDynamic) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType, c.Start, c.Mid, c.End)
		c.hash = &hash
	}
	return *c.hash
}

func (c *RuleIntBetweenDynamic) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
	if ok := facts.GetCacheInstance().TryGet(c.GetHash(), &isMatch); ok {
		return isMatch, nil
	}
	defer func() {
		facts.GetCacheInstance().SetCache(c.GetHash(), isMatch)
	}()
	param := facts.GetMap()

	// collecting values
	_params := [3]interface{}{c.Start, c.Mid, c.End}
	_values := [3]int{0, 0, 0}

	for i, _param := range _params {
		if err = common.ConvertInt().WithFromMap(param)(_param, &_values[i]); err != nil {
			return false, err
		}
	}

	return _values[0] <= _values[1] && _values[1] <= _values[2], nil
}
