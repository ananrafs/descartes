package rule_string

import (
	"strings"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleStringEqualFold struct {
	RuleType string `json:"type"`
	Field    string `json:"field"`
	Value    string `json:"value"`
}

func (c *RuleStringEqualFold) GetType() string {
	return "rules.string.equal_fold"
}

func (c *RuleStringEqualFold) New() rules.RulesItf {
	return new(RuleStringEqualFold)
}

func (c *RuleStringEqualFold) IsMatch(param map[string]interface{}) (isMatch bool, err error) {
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
