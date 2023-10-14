package rule_string

import (
	"strings"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleStringEqualFoldDynamic struct {
	RuleType string `json:"type"`
	Left     string `json:"left"`
	Right    string `json:"right"`
}

func (c *RuleStringEqualFoldDynamic) GetType() string {
	return "rules.string.equal_fold.dynamic"
}

func (c *RuleStringEqualFoldDynamic) New() rules.RulesItf {
	return new(RuleStringEqualFold)
}

func (c *RuleStringEqualFoldDynamic) IsMatch(param map[string]interface{}) (isMatch bool, err error) {
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
