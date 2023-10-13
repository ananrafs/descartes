package rule_string

import (
	"strings"

	"github.com/ananrafs/descartes/rules"
)

type RuleStringEqualFold struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    string `json:"value"`
}

func (c *RuleStringEqualFold) GetRules() string {
	return "rules.string.equal_fold"
}

func (c *RuleStringEqualFold) New() rules.RulesItf {
	return new(RuleStringEqualFold)
}

func (c *RuleStringEqualFold) IsMatch(param map[string]interface{}) (isMatch bool) {
	v, ok := param[c.Field]
	if !ok {
		return false
	}

	val, ok := v.(string)
	if !ok {
		return false
	}

	return strings.EqualFold(val, c.Value)
}
