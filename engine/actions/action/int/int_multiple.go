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
	for i, params := range c.Factors {
		val, numField := new(int), new(string)
		// first check if value was int
		if err = common.ConvertToInt(params, val); err == nil {
			if i == 0 {
				total = *val
			}
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

		total *= *val
	}

	param[c.Field] = total

	return
}
