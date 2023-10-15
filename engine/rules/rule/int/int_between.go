package rule_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleIntBetween struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Start    int    `json:"start"`
	End      int    `json:"end"`
	hash     *string
}

func (c *RuleIntBetween) GetType() string {
	return "rules.int.between"
}

func (c *RuleIntBetween) New() rules.RulesItf {
	return new(RuleIntBetween)
}

func (c *RuleIntBetween) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType, c.Field, c.Start, c.End)
		c.hash = &hash
	}
	return *c.hash
}

func (c *RuleIntBetween) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
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

	return c.Start <= *intv && *intv <= c.End, nil
}
