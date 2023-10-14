package group

import "github.com/ananrafs/descartes/engine/rules"

type ConditionalOr struct {
	ConditionalType string    `json:"type"`
	Rules           RuleGroup `json:"rules"`
}

func (c *ConditionalOr) GetType() string {
	return "rules.conditional.or"
}

func (c *ConditionalOr) New() rules.RulesItf {
	return new(ConditionalOr)
}

func (c *ConditionalOr) IsMatch(param map[string]interface{}) (isMatch bool, err error) {
	for _, rule := range c.Rules {
		isMatch, err = rule.IsMatch(param)
		if isMatch {
			return
		}
	}

	return
}
