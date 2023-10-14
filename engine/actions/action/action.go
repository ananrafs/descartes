package action

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/actions"
)

type Action map[string]interface{}

func (c *Action) GetType() string {
	return ""
}

func (c *Action) New() actions.ActionsItf {
	return new(Action)
}

func (c *Action) Do(param map[string]interface{}) (res interface{}, err error) {
	for i, params := range *c {
		paramsWithHandlebars := new(string)
		// check if its using handlebars
		if match := common.ParseFromHandlebars(params, paramsWithHandlebars); match {
			v, ok := param[*paramsWithHandlebars]
			if !ok {
				return nil, common.ErrorNotFoundOnMap(*paramsWithHandlebars)
			}

			(*c)[i] = v
		}

	}

	return *c, nil
}
