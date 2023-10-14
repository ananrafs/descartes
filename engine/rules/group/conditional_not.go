package group

import (
	"encoding/json"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/rules"
)

type ConditionalNot struct {
	RuleType string         `json:"type"`
	Rule     rules.RulesItf `json:"rule"`
}

func (c *ConditionalNot) GetType() string {
	return "rules.conditional.not"
}

func (c *ConditionalNot) New() rules.RulesItf {
	return new(ConditionalNot)
}

func (c *ConditionalNot) IsMatch(param map[string]interface{}) (isMatch bool, err error) {
	isMatch, _ = c.Rule.IsMatch(param)

	return !isMatch, nil
}

func (r *ConditionalNot) UnmarshalJSON(data []byte) (err error) {
	var m map[string]json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}

	for k, val := range m {
		var typeChecker common.TypeChecker

		switch k {
		case "type":
			if err := json.Unmarshal(val, &r.RuleType); err != nil {
				return err
			}
		case "rule":
			var instance rules.RulesItf

			if err := json.Unmarshal(val, &typeChecker); err != nil {
				return err
			}
			rule := rules.Get(typeChecker.Type)
			instance = rule.New()
			if err := json.Unmarshal(val, instance); err != nil {
				return err
			}
			r.Rule = instance
		}
	}

	return
}
