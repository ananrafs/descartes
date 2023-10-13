package rule_int

import "github.com/ananrafs/descartes/rules"

type RuleIntEqual struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    int    `json:"value"`
}

func (c *RuleIntEqual) GetRules() string {
	return "rules.int.equal"
}

func (c *RuleIntEqual) New() rules.RulesItf {
	return new(RuleIntEqual)
}

func (c *RuleIntEqual) IsMatch(param map[string]interface{}) (isMatch bool) {
	v, ok := param[c.Field]
	if !ok {
		return false
	}

	intf, ok := v.(float64)
	if !ok {
		return false
	}
	intv := int(intf)

	return intv == c.Value
}
