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
	dividedParams := [2]interface{}{c.Numerator, c.Denominator}
	dividedValues := [2]int{0, 0}

	for i, params := range dividedParams {
		val, numField := new(int), new(string)

		// first check if value was int
		if err = common.ConvertToInt(params, val); err == nil {
			dividedValues[i] = *val
			continue
		}

		// check if its using handlebars
		if match := common.ParseFromHandlebars(params, numField); !match {
			return nil, common.ErrorCasting(params)
		}

		v, ok := param[*numField]
		if !ok {
			return nil, common.ErrorNotFoundOnMap(*numField)
		}

		if err = common.ConvertToInt(v, val); err != nil {
			return nil, err
		}

		dividedValues[i] = *val
	}

	param[c.Field] = dividedValues[0] / dividedValues[1]

	return
}
