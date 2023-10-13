package group

import (
	"encoding/json"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/evaluators"
)

type EvaluatorGroup []evaluators.EvaluatorItf

func (eg *EvaluatorGroup) UnmarshalJSON(data []byte) (err error) {
	var m []json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}
	*eg = make(EvaluatorGroup, 0, len(m))
	var typeCheker common.TypeChecker
	for _, raw := range m {
		if err := json.Unmarshal(raw, &typeCheker); err != nil {
			return err
		}

		evals := evaluators.GetEvaluators(typeCheker.Type)
		newInstance := evals.New()
		err = json.Unmarshal(raw, newInstance)
		if err != nil {
			return err
		}
		*eg = append(*eg, newInstance)
	}
	return
}
