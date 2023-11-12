package group

import (
	"github.com/ananrafs/descartes/engine/evaluators"
	"github.com/ananrafs/descartes/engine/facts"
)

type FirstMatch struct {
	EvaluatorType string         `json:"type"`
	Evaluators    EvaluatorGroup `json:"evaluators"`
}

func (fm *FirstMatch) GetType() string {
	return "evaluator.group.first_match"
}

func NewFirstMatch() evaluators.EvaluatorsItf {
	o := new(FirstMatch)
	o.EvaluatorType = o.GetType()
	return o
}

func (fm *FirstMatch) Eval(facts facts.FactsItf) (res evaluators.EvalResult) {
	for _, eval := range fm.Evaluators {
		res = eval.Eval(facts)
		if res.IsMatch {
			return res
		}
	}
	return
}
