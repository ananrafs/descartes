package rule_string

import (
	"strings"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleStringEqualFold struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    string `json:"value"`
	hash     *string
}

func (c *RuleStringEqualFold) GetType() string {
	return "rules.string.equal_fold"
}

func (c *RuleStringEqualFold) New() rules.RulesItf {
	return new(RuleStringEqualFold)
}

func (c *RuleStringEqualFold) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType, c.Field, c.Value)
		c.hash = &hash
	}
	return *c.hash
}

func (c *RuleStringEqualFold) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
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

	return strings.EqualFold(val, c.Value), nil
}
