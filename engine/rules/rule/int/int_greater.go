package rule_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleIntGreater struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    int    `json:"value"`
	hash     *string
}

func (c *RuleIntGreater) GetType() string {
	return "rules.int.greater"
}

func (c *RuleIntGreater) New() rules.RulesItf {
	return new(RuleIntGreater)
}

func (c *RuleIntGreater) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType, c.Field, c.Value)
		c.hash = &hash
	}
	return *c.hash
}

func (c *RuleIntGreater) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
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

	intv := new(int)
	if err = common.ConvertToInt(v, intv); err != nil {
		return false, err
	}

	return *intv > c.Value, nil
}
