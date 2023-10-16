package action_float

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/actions"
	"github.com/ananrafs/descartes/engine/facts"
)

type Sum struct {
	Type    string        `json:"type"`
	Field   string        `json:"field"`
	Factors []interface{} `json:"factors"`
}

func (c *Sum) GetType() string {
	return "actions.float.sum"
}

func (c *Sum) New() actions.ActionsItf {
	return new(Sum)
}

func (c *Sum) Do(facts facts.FactsItf) (res interface{}, err error) {
	param := facts.GetMap()
	total := float64(0)
	for _, _param := range c.Factors {
		val := new(float64)
		if err = common.ConvertFloat().WithFromMap(param)(_param, val); err != nil {
			return false, err
		}
		total += *val
	}

	param[c.Field] = total

	return
}
