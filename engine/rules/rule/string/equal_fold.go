package rule_string

import (
	"strings"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type EqualFold struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    string `json:"value"`
	hash     *string
}

func (c *EqualFold) GetType() string {
	return "rules.string.equal_fold"
}

func (c *EqualFold) New() rules.RulesItf {
	o := new(EqualFold)
	o.RuleType = o.GetType()
	return o
}

func (c *EqualFold) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType, c.Field, c.Value)
		c.hash = &hash
	}
	return *c.hash
}

func (c *EqualFold) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
	if ok := facts.GetCacheInstance().TryGet(c.GetHash(), &isMatch); ok {
		return isMatch, nil
	}
	defer func() {
		facts.GetCacheInstance().SetCache(c.GetHash(), isMatch)
	}()
	param := facts.GetMap()

	v, err := common.LookUpRecursiveMap(param, c.Field)
	if err != nil {
		return false, common.ErrorNotFoundOnMap(c.Field)
	}

	val, ok := v.(string)
	if !ok {
		return false, common.ErrorCasting(v)
	}

	return strings.EqualFold(val, c.Value), nil
}
