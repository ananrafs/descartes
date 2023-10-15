package action_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/actions"
	"github.com/ananrafs/descartes/engine/facts"
)

// 4 % 5.
// 4 is Dividend.
// 5 is Divisor
type ActionIntMod struct {
	Type     string      `json:"type"`
	Field    string      `json:"field"`
	Dividend interface{} `json:"dividend"`
	Divisor  interface{} `json:"divisor"`
}

func (c *ActionIntMod) GetType() string {
	return "actions.int.mod"
}

func (c *ActionIntMod) New() actions.ActionsItf {
	return new(ActionIntMod)
}

func (c *ActionIntMod) Do(facts facts.FactsItf) (res interface{}, err error) {
	// collecting values
	_params := [2]interface{}{c.Dividend, c.Divisor}
	_values := [2]int{0, 0}
	param := facts.GetMap()

	for i, _param := range _params {
		if err = common.ConvertInt().WithFromMap(param)(_param, &_values[i]); err != nil {
			return false, err
		}
	}

	param[c.Field] = _values[0] % _values[1]

	return
}
