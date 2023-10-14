package rule_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleIntGreaterDynamic struct {
	RuleType string `json:"type"`
	Left     string `json:"left"`
	Right    int    `json:"right"`
}

func (c *RuleIntGreaterDynamic) GetType() string {
	return "rules.int.greater.dynamic"
}

func (c *RuleIntGreaterDynamic) New() rules.RulesItf {
	return new(RuleIntGreaterDynamic)
}

func (c *RuleIntGreaterDynamic) IsMatch(param map[string]interface{}) (isMatch bool, err error) {
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
