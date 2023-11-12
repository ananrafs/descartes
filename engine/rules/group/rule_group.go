package group

import (
	"encoding/json"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/rules"
)

type RuleGroup []rules.RulesItf

func (r *RuleGroup) UnmarshalJSON(data []byte) (err error) {
	var m []json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}
	*r = make(RuleGroup, 0, len(m))
	for _, raw := range m {
		var typeChecker common.TypeChecker

		if err := json.Unmarshal(raw, &typeChecker); err != nil {
			return err
		}

		newInstance := rules.Get(typeChecker.Type)
		err = json.Unmarshal(raw, newInstance)
		if err != nil {
			return err
		}

		*r = append(*r, newInstance)
	}
	return
}
