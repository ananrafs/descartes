package group

import (
	"github.com/ananrafs/descartes/engine/actions"
)

type ActionGroup struct {
	Type    string       `json:"type"`
	Actions ActionsGroup `json:"actions"`
}

func (c *ActionGroup) GetType() string {
	return "actions.group"
}

func (c *ActionGroup) New() actions.ActionsItf {
	return new(ActionGroup)
}

func (c *ActionGroup) Do(param map[string]interface{}) (res interface{}, err error) {
	for _, action := range c.Actions {
		res, err = action.Do(param)
		if err != nil {
			return nil, err
		}
	}

	return
}
