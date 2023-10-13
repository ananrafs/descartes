package rule_int

import "github.com/ananrafs/descartes/rules"

type RuleIntGreater struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    int    `json:"value"`
}

func (c *RuleIntGreater) GetRules() string {
	return "rules.int.greater"
}

func (c *RuleIntGreater) New() rules.RulesItf {
	return new(RuleIntGreater)
}

func (c *RuleIntGreater) IsMatch(param map[string]interface{}) (isMatch bool) {
	v, ok := param[c.Field]
	if !ok {
		return false
	}

	intf, ok := v.(float64)
	if !ok {
		return false
	}
	intv := int(intf)

	return intv > c.Value
}
