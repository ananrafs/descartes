package group

import (
	"github.com/ananrafs/descartes/engine/evaluators"
)

type FirstMatch struct {
	EvalType   string         `json:"type"`
	Evaluators EvaluatorGroup `json:"evaluators"`
}

func (fm *FirstMatch) GetType() string {
	return "evaluator.group.first_match"
}

func (fm *FirstMatch) New() evaluators.EvaluatorsItf {
	return new(FirstMatch)
}

func (fm *FirstMatch) Eval(param map[string]interface{}) (res evaluators.EvalResult) {
	for _, eval := range fm.Evaluators {
		res = eval.Eval(param)
		if res.IsMatch || res.Error != nil {
			return res
		}
	}
	return
}
