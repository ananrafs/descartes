package group

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type ConditionalOr struct {
	ConditionalType string    `json:"type"`
	Rules           RuleGroup `json:"rules"`
	hash            *string
}

func (c *ConditionalOr) GetType() string {
	return "rules.conditional.or"
}

func (c *ConditionalOr) New() rules.RulesItf {
	o := new(ConditionalOr)
	o.ConditionalType = o.GetType()
	return o
}

func (c *ConditionalOr) GetHash() string {
	for c.hash == nil {
		hashs := make([]interface{}, 0, len(c.Rules)+1)
		hashs = append(hashs, c.ConditionalType)
		for _, rule := range c.Rules {
			hashs = append(hashs, rule.GetHash())
		}

		hash := common.CreateHash(hashs...)
		c.hash = &hash
	}
	return *c.hash
}

func (c *ConditionalOr) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
	if ok := facts.GetCacheInstance().TryGet(c.GetHash(), &isMatch); ok {
		return isMatch, nil
	}
	defer func() {
		facts.GetCacheInstance().SetCache(c.GetHash(), isMatch)
	}()

	for _, rule := range c.Rules {
		isMatch, err = rule.IsMatch(facts)
		if isMatch {
			return
		}
	}

	return
}
