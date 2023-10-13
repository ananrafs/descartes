package group

import (
	"github.com/ananrafs/descartes/rules"
)

type ConditionalAnd struct {
	ConditionalType string    `json:"type"`
	Rules           RuleGroup `json:"rules"`
}

func (c *ConditionalAnd) GetRules() string {
	return "rules.conditional.and"
}

func (c *ConditionalAnd) New() rules.RulesItf {
	return new(ConditionalAnd)
}

func (c *ConditionalAnd) IsMatch(param map[string]interface{}) (isMatch bool) {
	isMatch = true
	for _, rule := range c.Rules {
		isMatch = rule.IsMatch(param)
		if !isMatch {
			return
		}
	}
	return
}
