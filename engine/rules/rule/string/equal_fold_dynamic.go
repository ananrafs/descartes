package rule_string

import (
	"fmt"
	"strings"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type EqualFoldDynamic struct {
	RuleType string `json:"type"`
	Left     string `json:"left"`
	Right    string `json:"right"`
	hash     *string
}

func (c *EqualFoldDynamic) GetType() string {
	return "rules.string.equal_fold.dynamic"
}

func (c *EqualFoldDynamic) New() rules.RulesItf {
	o := new(EqualFoldDynamic)
	o.RuleType = o.GetType()
	return o
}

func (c *EqualFoldDynamic) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType, c.Left, c.Right)
		c.hash = &hash
	}
	return *c.hash
}

func (c *EqualFoldDynamic) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
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

	return strings.EqualFold(_values[0], _values[1]), nil
}
