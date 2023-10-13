package rule_string

import "github.com/ananrafs/descartes/rules"

type RuleStringEqual struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    string `json:"value"`
}

func (c *RuleStringEqual) GetRules() string {
	return "rules.string.equal"
}

func (c *RuleStringEqual) New() rules.RulesItf {
	return new(RuleStringEqual)
}

func (c *RuleStringEqual) IsMatch(param map[string]interface{}) (isMatch bool) {
	v, ok := param[c.Field]
	if !ok {
		return false
	}

	val, ok := v.(string)
	if !ok {
		return false
	}

	return val == c.Value
}
