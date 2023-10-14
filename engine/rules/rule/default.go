package rule

import (
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleDefault struct {
	RuleType string `json:"type"`
}

func (c *RuleDefault) GetType() string {
	return "rules.default"
}

func (c *RuleDefault) New() rules.RulesItf {
	return new(RuleDefault)
}

func (c *RuleDefault) IsMatch(param map[string]interface{}) (isMatch bool, err error) {
	return true, nil
}
