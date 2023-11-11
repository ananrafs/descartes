package group

import (
	"github.com/ananrafs/descartes/engine/actions"
	"github.com/ananrafs/descartes/engine/facts"
)

type ActionGroup struct {
	Type    string       `json:"type"`
	Actions ActionsGroup `json:"actions"`
}

func (c *ActionGroup) GetType() string {
	return "actions.group"
}

func (c *ActionGroup) New() actions.ActionsItf {
	o := new(ActionGroup)
	o.Type = o.GetType()
	return o
}

func (c *ActionGroup) Do(facts facts.FactsItf) (res interface{}, err error) {
	for _, action := range c.Actions {
		res, err = action.Do(facts)
		// if err != nil {
		// 	return nil, err
		// }
	}

	return
}
