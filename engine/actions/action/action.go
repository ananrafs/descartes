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
		paramsWithTemplate := new(string)
		// check if its using template
		if match := common.ParseFromMustacheTemplate(params, paramsWithTemplate); match {
			v, ok := param[*paramsWithTemplate]
			if !ok {
				return nil, common.ErrorNotFoundOnMap(*paramsWithTemplate)
			}

			(*c)[i] = v
		}

	}

	return *c, nil
}
