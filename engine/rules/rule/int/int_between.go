package rule_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleIntBetween struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Start    int    `json:"start"`
	End      int    `json:"end"`
}

func (c *RuleIntBetween) GetType() string {
	return "rules.int.between"
}

func (c *RuleIntBetween) New() rules.RulesItf {
	return new(RuleIntBetween)
}

func (c *RuleIntBetween) IsMatch(param map[string]interface{}) (isMatch bool, err error) {
	v, ok := param[c.Field]
	if !ok {
		return false, common.ErrorNotFoundOnMap(c.Field)
	}

	intv := new(int)
	if err = common.ConvertToInt(v, intv); err != nil {
		return false, err
	}

	return c.Start <= *intv && *intv <= c.End, nil
}
