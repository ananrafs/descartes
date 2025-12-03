package rule_string

import (
	"fmt"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type EqualWildcardDynamic struct {
	RuleType string `json:"type"`
	Left     string `json:"left"`
	Right    string `json:"right"`
	hash     *string
}

func (c *EqualWildcardDynamic) GetType() string {
	return "rules.string.equal.wildcard.dynamic"
}

func NewEqualWildcardDynamic() rules.RulesItf {
	o := new(EqualWildcardDynamic)
	o.RuleType = o.GetType()
	return o
}

func (c *EqualWildcardDynamic) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType, c.Left, c.Right)
		c.hash = &hash
	}
	return *c.hash
}

func (c *EqualWildcardDynamic) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
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
		key := ""
		var _field interface{}
		if match := common.DeepTemplateEvaluateFromMap(param, _param, &_field); match {
			key = fmt.Sprintf("%v", _field)
			var ok bool
			_param, ok = param[key]
			if !ok {
				return false, common.ErrorNotFoundOnMap(key)
			}
		}

		val, ok := _param.(string)
		if !ok {
			return false, common.ErrorCasting(_param)
		}
		_values[i] = val
	}

	// Check if either value is wildcard
	if _values[0] == "ALL" || _values[0] == "*" || _values[1] == "ALL" || _values[1] == "*" {
		return true, nil
	}

	return _values[0] == _values[1], nil
}
