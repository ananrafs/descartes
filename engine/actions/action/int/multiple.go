package action_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/actions"
	"github.com/ananrafs/descartes/engine/facts"
)

type Multiple struct {
	Type    string        `json:"type"`
	Field   string        `json:"field"`
	Factors []interface{} `json:"factors"`
}

func (c *Multiple) GetType() string {
	return "actions.int.multiple"
}

func (c *Multiple) New() actions.ActionsItf {
	o := new(Multiple)
	o.Type = o.GetType()
	o.Factors = make([]interface{}, 0)
	return o
}

func (c *Multiple) Do(facts facts.FactsItf) (res interface{}, err error) {
	param := facts.GetMap()
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
