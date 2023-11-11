package group

import (
	"encoding/json"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type ConditionalNot struct {
	ConditionalType string         `json:"type"`
	Rule            rules.RulesItf `json:"rule"`
	hash            *string
}

func (c *ConditionalNot) GetType() string {
	return "rules.conditional.not"
}

func (c *ConditionalNot) New() rules.RulesItf {
	o := new(ConditionalNot)
	o.ConditionalType = o.GetType()
	return o
}

func (c *ConditionalNot) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.ConditionalType, c.Rule.GetHash())
		c.hash = &hash
	}
	return *c.hash

}

func (c *ConditionalNot) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
	if ok := facts.GetCacheInstance().TryGet(c.GetHash(), &isMatch); ok {
		return isMatch, nil
	}
	defer func() {
		facts.GetCacheInstance().SetCache(c.GetHash(), isMatch)
	}()

	ok, _ := c.Rule.IsMatch(facts)

	return !ok, nil
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
			if err := json.Unmarshal(val, &r.ConditionalType); err != nil {
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
