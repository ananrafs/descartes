package action_float

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
	return "actions.float.multiple"
}

func NewMultiple() actions.ActionsItf {
	o := new(Multiple)
	o.Type = o.GetType()
	o.Factors = make([]interface{}, 0)
	return o
}

func (c *Multiple) Do(facts facts.FactsItf) (res interface{}, err error) {
	param := facts.GetMap()
	total := float64(0)
	for i, _param := range c.Factors {
		val := new(float64)
		if err = common.Convert[float64]().WithFromMap(param)(_param, val); err != nil {
			return false, err
		}
		if i == 0 {
			total = *val
		} else {
			total *= *val
		}
	}
	common.SetMap(param, c.Field, total)

	return
}
