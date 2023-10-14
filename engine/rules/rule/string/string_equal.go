package rule_string

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleStringEqual struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    string `json:"value"`
}

func (c *RuleStringEqual) GetType() string {
	return "rules.string.equal"
}

func (c *RuleStringEqual) New() rules.RulesItf {
	return new(RuleStringEqual)
}

func (c *RuleStringEqual) IsMatch(param map[string]interface{}) (isMatch bool, err error) {
	v, ok := param[c.Field]
	if !ok {
		return false, common.ErrorNotFoundOnMap(c.Field)
	}

	val, ok := v.(string)
	if !ok {
		return false, common.ErrorCasting(v)
	}

	return val == c.Value, nil
}
