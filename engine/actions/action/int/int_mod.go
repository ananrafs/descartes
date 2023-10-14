package action_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/actions"
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

func (c *ActionIntMod) Do(param map[string]interface{}) (res interface{}, err error) {
	// collecting values
	modParams := [2]interface{}{c.Dividend, c.Divisor}
	modValues := [2]int{0, 0}

	for i, params := range modParams {
		val, numField := new(int), new(string)
		// first check if value was int
		if err = common.ConvertToInt(params, val); err == nil {
			modValues[i] = *val
			continue
		}

		// check if its using template
		if match := common.ParseFromMustacheTemplate(params, numField); !match {
			return nil, common.ErrorCasting(params)
		}

		v, ok := param[*numField]
		if !ok {
			return nil, common.ErrorNotFoundOnMap(*numField)
		}

		if err = common.ConvertToInt(v, val); err != nil {
			return nil, err
		}
		modValues[i] = *val
	}

	param[c.Field] = modValues[0] % modValues[1]

	return
}
