package group

import (
	"github.com/ananrafs/descartes/engine/rules"
)

type ConditionalAnd struct {
	ConditionalType string    `json:"type"`
	Rules           RuleGroup `json:"rules"`
}

func (c *ConditionalAnd) GetType() string {
	return "rules.conditional.and"
}

func (c *ConditionalAnd) New() rules.RulesItf {
	return new(ConditionalAnd)
}

func (c *ConditionalAnd) IsMatch(param map[string]interface{}) (isMatch bool, err error) {
	isMatch = true
	for _, rule := range c.Rules {
		isMatch, err = rule.IsMatch(param)
		if !isMatch {
			return false, nil
		}
	}
	return
}
