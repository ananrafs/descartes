package rule_int

import "github.com/ananrafs/descartes/rules"

type RuleIntBetween struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Start    int    `json:"start"`
	End      int    `json:"end"`
}

func (c *RuleIntBetween) GetRules() string {
	return "rules.int.between"
}

func (c *RuleIntBetween) New() rules.RulesItf {
	return new(RuleIntBetween)
}

func (c *RuleIntBetween) IsMatch(param map[string]interface{}) (isMatch bool) {
	v, ok := param[c.Field]
	if !ok {
		return false
	}
	intf, ok := v.(float64)
	if !ok {
		return false
	}
	intv := int(intf)

	return c.Start <= intv && intv <= c.End
}
