package action_int

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/actions"
	"github.com/ananrafs/descartes/engine/facts"
)

type ActionIntSum struct {
	Type    string        `json:"type"`
	Field   string        `json:"field"`
	Factors []interface{} `json:"factors"`
}

func (c *ActionIntSum) GetType() string {
	return "actions.int.sum"
}

func (c *ActionIntSum) New() actions.ActionsItf {
	return new(ActionIntSum)
}

func (c *ActionIntSum) Do(facts facts.FactsItf) (res interface{}, err error) {
	param := facts.GetMap()
	total := 0
	for _, _param := range c.Factors {
		val := new(int)
		if err = common.ConvertInt().WithFromMap(param)(_param, val); err != nil {
			return false, err
		}
		total += *val
	}

	param[c.Field] = total

	return
}
