package rule_bool

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type Bool struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    bool   `json:"value"`
	hash     *string
}

func (c *Bool) GetType() string {
	return "rules.bool"
}

func (c *Bool) New() rules.RulesItf {
	return new(Bool)
}

func (c *Bool) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType, c.Field, c.Value)
		c.hash = &hash
	}
	return *c.hash
}

func (c *Bool) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
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

	isTrue := new(bool)
	if err = common.ConvertToBool(v, isTrue); err != nil {
		return false, err
	}

	return *isTrue == c.Value, nil
}
