package rule

import (
	"github.com/ananrafs/descartes/cache"
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleDefault struct {
	RuleType string `json:"type"`
	hash     *string
}

func (c *RuleDefault) GetType() string {
	return "rules.default"
}

func (c *RuleDefault) New() rules.RulesItf {
	return new(RuleDefault)
}

func (c *RuleDefault) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType)
		c.hash = &hash
	}
	return *c.hash
}

// no need cache default rules,
// so doesnt need cache
func (c *RuleDefault) SetCache(_cache cache.CacheItf) {
	return
}

func (c *RuleDefault) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
	return true, nil
}
