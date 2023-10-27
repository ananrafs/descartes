package group

import (
	"github.com/ananrafs/descartes/engine/evaluators"
	"github.com/ananrafs/descartes/engine/facts"
)

type MultiMatchOrdered struct {
	EvalType string `json:"type"`

	// maximum matched evaluation allowed
	MaxMatch int `json:"max"`

	// reentrance = true,
	// means rules matched will reevaluate.
	// set false, means matched rules will taked out from evaluate group
	Reentrance bool `json:"reentrance"`

	// determine if evalresult should be merged
	Merging bool `json:"merging"`

	Evaluators EvaluatorGroup `json:"evaluators"`
}

func (fm *MultiMatchOrdered) GetType() string {
	return "evaluator.group.multi_match_ordered"
}

func (fm *MultiMatchOrdered) New() evaluators.EvaluatorsItf {
	return new(MultiMatchOrdered)
}

func (fm MultiMatchOrdered) Eval(facts facts.FactsItf) (res evaluators.EvalResult) {

	deduct := func(instance *int) (deducted bool) {
		(*instance)--
		return true
	}

	mapEvaluators := make(map[int]evaluators.EvaluatorsItf)
	evaluatedMap := make(map[int]bool)

	for index, eval := range fm.Evaluators {
		mapEvaluators[index] = eval
	}

	var response evaluators.EvalResult
	for fm.MaxMatch > 0 {

		var deducted bool
		for index := 0; index < len(mapEvaluators); index++ {
			if _, ok := evaluatedMap[index]; ok && fm.Reentrance {
				continue
			}
			eval := mapEvaluators[index]

			response = eval.Eval(facts)
			if response.IsMatch {
				if !fm.Reentrance {
					evaluatedMap[index] = true
				}
				deducted = deduct(&fm.MaxMatch)
				if fm.Merging {
					res.Merge(response)
					continue
				}
				res = response
			}
		}

		// break loop if there's no match evaluation in one cycle
		if !deducted {
			break
		}
	}

	return

}
