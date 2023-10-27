package group

import (
	"encoding/json"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/evaluators"
)

type EvaluatorGroup []evaluators.EvaluatorsItf

func (eg *EvaluatorGroup) UnmarshalJSON(data []byte) (err error) {
	var m []json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}
	*eg = make(EvaluatorGroup, 0, len(m))

	for _, raw := range m {
		var typeChecker common.TypeChecker
		if err := json.Unmarshal(raw, &typeChecker); err != nil {
			return err
		}
		evals := evaluators.Get(typeChecker.Type)
		newInstance := evals.New()
		err = json.Unmarshal(raw, newInstance)
		if err != nil {
			return err
		}
		*eg = append(*eg, newInstance)
	}
	return
}
