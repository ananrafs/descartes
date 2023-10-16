package action_float

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/actions"
	"github.com/ananrafs/descartes/engine/facts"
)

type Substract struct {
	Type    string        `json:"type"`
	Field   string        `json:"field"`
	Factors []interface{} `json:"factors"`
}

func (c *Substract) GetType() string {
	return "actions.float.substract"
}

func (c *Substract) New() actions.ActionsItf {
	return new(Substract)
}

func (c *Substract) Do(facts facts.FactsItf) (res interface{}, err error) {
	param := facts.GetMap()
	total := float64(0)
	for i, _param := range c.Factors {
		val := new(float64)
		if err = common.ConvertFloat().WithFromMap(param)(_param, val); err != nil {
			return false, err
		}

		if i == 0 {
			total = *val
		} else {
			total -= *val
		}
	}

	param[c.Field] = total

	return
}
