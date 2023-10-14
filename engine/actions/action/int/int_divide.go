package action_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/actions"
)

type ActionIntDivide struct {
	Type        string      `json:"type"`
	Field       string      `json:"field"`
	Numerator   interface{} `json:"numerator"`
	Denominator interface{} `json:"denominator"`
}

func (c *ActionIntDivide) GetType() string {
	return "actions.int.divide"
}

func (c *ActionIntDivide) New() actions.ActionsItf {
	return new(ActionIntDivide)
}

func (c *ActionIntDivide) Do(param map[string]interface{}) (res interface{}, err error) {
	// collecting values
	_params := [2]interface{}{c.Numerator, c.Denominator}
	_values := [2]int{0, 0}

	for i, _param := range _params {
		if err = common.ConvertInt().WithFromMap(param)(_param, &_values[i]); err != nil {
			return false, err
		}
	}

	param[c.Field] = _values[0] / _values[1]

	return
}
