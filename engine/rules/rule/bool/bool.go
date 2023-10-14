package rule_bool

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleBool struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    bool   `json:"value"`
}

func (c *RuleBool) GetType() string {
	return "rules.bool"
}

func (c *RuleBool) New() rules.RulesItf {
	return new(RuleBool)
}

func (c *RuleBool) IsMatch(param map[string]interface{}) (isMatch bool, err error) {
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
