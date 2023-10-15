package rule_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleIntGreaterDynamic struct {
	RuleType string `json:"type"`
	Left     string `json:"left"`
	Right    int    `json:"right"`
	hash     *string
}

func (c *RuleIntGreaterDynamic) GetType() string {
	return "rules.int.greater.dynamic"
}

func (c *RuleIntGreaterDynamic) New() rules.RulesItf {
	return new(RuleIntGreaterDynamic)
}

func (c *RuleIntGreaterDynamic) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType, c.Left, c.Right)
		c.hash = &hash
	}
	return *c.hash
}

func (c *RuleIntGreaterDynamic) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
	if ok := facts.GetCacheInstance().TryGet(c.GetHash(), &isMatch); ok {
		return isMatch, nil
	}
	defer func() {
		facts.GetCacheInstance().SetCache(c.GetHash(), isMatch)
	}()
	param := facts.GetMap()

	// collecting values
	_params := [2]interface{}{c.Left, c.Right}
	_values := [2]int{0, 0}

	for i, _param := range _params {
		if err = common.ConvertInt().WithFromMap(param)(_param, &_values[i]); err != nil {
			return false, err
		}
	}
	return _values[0] > _values[1], nil
}
