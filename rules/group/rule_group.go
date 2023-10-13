package group

import (
	"encoding/json"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/rules"
)

type RuleGroup []rules.RulesItf

func (r *RuleGroup) UnmarshalJSON(data []byte) (err error) {
	var m []json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}
	*r = make(RuleGroup, 0, len(m))
	var typeCheker common.TypeChecker
	for _, raw := range m {
		if err := json.Unmarshal(raw, &typeCheker); err != nil {
			return err
		}

		rule := rules.GetRules(typeCheker.Type)
		newInstance := rule.New()
		err = json.Unmarshal(raw, newInstance)
		if err != nil {
			return err
		}

		*r = append(*r, newInstance)
	}
	return
}
