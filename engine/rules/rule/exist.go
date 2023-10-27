package rule

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type Exist struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	hash     *string
}

func (c *Exist) GetType() string {
	return "rules.exist"
}

func (c *Exist) New() rules.RulesItf {
	return new(Exist)
}

func (c *Exist) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType)
		c.hash = &hash
	}
	return *c.hash
}

func (c *Exist) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
	if ok := facts.GetCacheInstance().TryGet(c.GetHash(), &isMatch); ok {
		return isMatch, nil
	}
	defer func() {
		facts.GetCacheInstance().SetCache(c.GetHash(), isMatch)
	}()
	param := facts.GetMap()

	_, ok := param[c.Field]
	return ok, nil
}
