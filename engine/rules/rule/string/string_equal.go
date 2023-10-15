package rule_string

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleStringEqual struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    string `json:"value"`
	hash     *string
}

func (c *RuleStringEqual) GetType() string {
	return "rules.string.equal"
}

func (c *RuleStringEqual) New() rules.RulesItf {
	return new(RuleStringEqual)
}

func (c *RuleStringEqual) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType, c.Field, c.Value)
		c.hash = &hash
	}
	return *c.hash
}

func (c *RuleStringEqual) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
	if ok := facts.GetCacheInstance().TryGet(c.GetHash(), &isMatch); ok {
		return isMatch, nil
	}
	defer func() {
		facts.GetCacheInstance().SetCache(c.GetHash(), isMatch)
	}()
	param := facts.GetMap()

	v, ok := param[c.Field]
	if !ok {
		return false, common.ErrorNotFoundOnMap(c.Field)
	}

	val, ok := v.(string)
	if !ok {
		return false, common.ErrorCasting(v)
	}

	return val == c.Value, nil
}
