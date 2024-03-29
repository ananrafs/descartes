package rule_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type LesserEqualDynamic struct {
	RuleType string      `json:"type"`
	Left     interface{} `json:"left"`
	Right    interface{} `json:"right"`
	hash     *string
}

func (c *LesserEqualDynamic) GetType() string {
	return "rules.int.lesser_equal.dynamic"
}

func NewLesserEqualDynamic() rules.RulesItf {
	o := new(LesserEqualDynamic)
	o.RuleType = o.GetType()
	return o
}

func (c *LesserEqualDynamic) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType, c.Left, c.Right)
		c.hash = &hash
	}
	return *c.hash
}

func (c *LesserEqualDynamic) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
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
		if err = common.Convert[int]().WithFromMap(param)(_param, &_values[i]); err != nil {
			return false, err
		}
	}

	return _values[0] <= _values[1], nil
}
