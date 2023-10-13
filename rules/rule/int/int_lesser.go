package rule_int

import "github.com/ananrafs/descartes/rules"

type RuleIntLesser struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    int    `json:"value"`
}

func (c *RuleIntLesser) GetRules() string {
	return "rules.int.lesser"
}

func (c *RuleIntLesser) New() rules.RulesItf {
	return new(RuleIntLesser)
}

func (c *RuleIntLesser) IsMatch(param map[string]interface{}) (isMatch bool) {
	v, ok := param[c.Field]
	if !ok {
		return false
	}

	intf, ok := v.(float64)
	if !ok {
		return false
	}
	intv := int(intf)

	return intv < c.Value
}
