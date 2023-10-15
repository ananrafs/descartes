package rule_string

import (
	"strings"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleStringEqualFoldDynamic struct {
	RuleType string `json:"type"`
	Left     string `json:"left"`
	Right    string `json:"right"`
	hash     *string
}

func (c *RuleStringEqualFoldDynamic) GetType() string {
	return "rules.string.equal_fold.dynamic"
}

func (c *RuleStringEqualFoldDynamic) New() rules.RulesItf {
	return new(RuleStringEqualFoldDynamic)
}

func (c *RuleStringEqualFoldDynamic) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType, c.Left, c.Right)
		c.hash = &hash
	}
	return *c.hash
}

func (c *RuleStringEqualFoldDynamic) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
	if ok := facts.GetCacheInstance().TryGet(c.GetHash(), &isMatch); ok {
		return isMatch, nil
	}
	defer func() {
		facts.GetCacheInstance().SetCache(c.GetHash(), isMatch)
	}()
	param := facts.GetMap()

	// collecting values
	_params := [2]interface{}{c.Left, c.Right}
	_values := [2]string{"", ""}

	for i, _param := range _params {
		_field := new(string)
		if match := common.ParseFromMustacheTemplate(_param, _field); match {
			var ok bool
			_param, ok = param[*_field]
			if !ok {
				return false, common.ErrorNotFoundOnMap(*_field)
			}
		}

		val, ok := _param.(string)
		if !ok {
			return false, common.ErrorCasting(_param)
		}
		_values[i] = val
	}

	return strings.EqualFold(_values[0], _values[1]), nil
}
