package rule_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleIntGreater struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    int    `json:"value"`
}

func (c *RuleIntGreater) GetType() string {
	return "rules.int.greater"
}

func (c *RuleIntGreater) New() rules.RulesItf {
	return new(RuleIntGreater)
}

func (c *RuleIntGreater) IsMatch(param map[string]interface{}) (isMatch bool, err error) {
	v, ok := param[c.Field]
	if !ok {
		return false, common.ErrorNotFoundOnMap(c.Field)
	}

	intv := new(int)
	if err = common.ConvertToInt(v, intv); err != nil {
		return false, err
	}

	return *intv > c.Value, nil
}
