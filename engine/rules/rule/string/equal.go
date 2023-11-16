package rule_string

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type Equal struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    string `json:"value"`
	hash     *string
}

func (c *Equal) GetType() string {
	return "rules.string.equal"
}

func NewEqual() rules.RulesItf {
	o := new(Equal)
	o.RuleType = o.GetType()
	return o
}

func (c *Equal) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType, c.Field, c.Value)
		c.hash = &hash
	}
	return *c.hash
}

func (c *Equal) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
	if ok := facts.GetCacheInstance().TryGet(c.GetHash(), &isMatch); ok {
		return isMatch, nil
	}
	defer func() {
		facts.GetCacheInstance().SetCache(c.GetHash(), isMatch)
	}()
	param := facts.GetMap()

	v, err := common.LookUpMap(param, c.Field)
	if err != nil {
		return false, common.ErrorNotFoundOnMap(c.Field)
	}

	val, ok := v.(string)
	if !ok {
		return false, common.ErrorCasting(v)
	}

	return val == c.Value, nil
}
