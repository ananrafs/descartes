package rule_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type LesserEqual struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    int    `json:"value"`
	hash     *string
}

func (c *LesserEqual) GetType() string {
	return "rules.int.lesser_equal"
}

func NewLesserEqual() rules.RulesItf {
	o := new(LesserEqual)
	o.RuleType = o.GetType()
	return o
}

func (c *LesserEqual) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType, c.Field, c.Value)
		c.hash = &hash
	}
	return *c.hash
}

func (c *LesserEqual) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
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

	intv := new(int)
	if err = common.Convert[int]()(v, intv); err != nil {
		return false, err
	}

	return *intv <= c.Value, nil
}
