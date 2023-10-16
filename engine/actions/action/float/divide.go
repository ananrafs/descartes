package action_float

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/actions"
	"github.com/ananrafs/descartes/engine/facts"
)

type Divide struct {
	Type        string      `json:"type"`
	Field       string      `json:"field"`
	Numerator   interface{} `json:"numerator"`
	Denominator interface{} `json:"denominator"`
}

func (c *Divide) GetType() string {
	return "actions.float.divide"
}

func (c *Divide) New() actions.ActionsItf {
	return new(Divide)
}

func (c *Divide) Do(facts facts.FactsItf) (res interface{}, err error) {
	// collecting values
	_params := [2]interface{}{c.Numerator, c.Denominator}
	_values := [2]float64{0, 0}
	param := facts.GetMap()

	for i, _param := range _params {
		if err = common.ConvertFloat().WithFromMap(param)(_param, &_values[i]); err != nil {
			return false, err
		}
	}

	param[c.Field] = _values[0] / _values[1]

	return
}
