package group

import (
	"encoding/json"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/actions"
)

type ActionsGroup []actions.ActionsItf

func (ag *ActionsGroup) UnmarshalJSON(data []byte) (err error) {
	var m []json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}
	*ag = make(ActionsGroup, 0, len(m))

	for _, raw := range m {
		var typeChecker common.TypeChecker
		if err := json.Unmarshal(raw, &typeChecker); err != nil {
			return err
		}

		action := actions.Get(typeChecker.Type)
		newInstance := action.New()
		err = json.Unmarshal(raw, newInstance)
		if err != nil {
			return err
		}
		*ag = append(*ag, newInstance)
	}
	return
}
