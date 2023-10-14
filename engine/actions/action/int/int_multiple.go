package action_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/actions"
)

type ActionIntMultiple struct {
	Type    string        `json:"type"`
	Field   string        `json:"field"`
	Factors []interface{} `json:"factors"`
}

func (c *ActionIntMultiple) GetType() string {
	return "actions.int.multiple"
}

func (c *ActionIntMultiple) New() actions.ActionsItf {
	return new(ActionIntMultiple)
}

func (c *ActionIntMultiple) Do(param map[string]interface{}) (res interface{}, err error) {
	total := 0
	for i, _param := range c.Factors {
		val := new(int)
		if err = common.ConvertInt().WithFromMap(param)(_param, val); err != nil {
			return false, err
		}
		if i == 0 {
			total = *val
		} else {
			total *= *val
		}
	}
	param[c.Field] = total

	return
}
