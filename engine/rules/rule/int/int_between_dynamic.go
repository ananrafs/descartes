package rule_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleIntBetweenDynamic struct {
	RuleType string      `json:"type"`
	Mid      interface{} `json:"mid"`
	Start    interface{} `json:"start"`
	End      interface{} `json:"end"`
}

func (c *RuleIntBetweenDynamic) GetType() string {
	return "rules.int.between.dynamic"
}

func (c *RuleIntBetweenDynamic) New() rules.RulesItf {
	return new(RuleIntBetween)
}

func (c *RuleIntBetweenDynamic) IsMatch(param map[string]interface{}) (isMatch bool, err error) {
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
