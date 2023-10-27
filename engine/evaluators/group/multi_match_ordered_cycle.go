package group

import (
	"github.com/ananrafs/descartes/engine/evaluators"
	"github.com/ananrafs/descartes/engine/facts"
)

type MultiMatchOrderedCycle struct {
	EvalType string `json:"type"`

	// maximum matched evaluation allowed
	MaxCycle int `json:"max"`

	// reentrance = true,
	// means rules matched will reevaluate.
	// set false, means matched rules will taked out from evaluate group
	Reentrance bool `json:"reentrance"`

	// determine if evalresult should be merged
	Merging bool `json:"merging"`

	Evaluators EvaluatorGroup `json:"evaluators"`
}

func (fm *MultiMatchOrderedCycle) GetType() string {
	return "evaluator.group.multi_match_ordered_cycle"
}

func (fm *MultiMatchOrderedCycle) New() evaluators.EvaluatorsItf {
	return new(MultiMatchOrderedCycle)
}

func (fm MultiMatchOrderedCycle) Eval(facts facts.FactsItf) (res evaluators.EvalResult) {
	mapEvaluators := make(map[int]evaluators.EvaluatorsItf)
	evaluatedMap := make(map[int]bool)

	for index, eval := range fm.Evaluators {
		mapEvaluators[index] = eval
	}

	var response evaluators.EvalResult
	for fm.MaxCycle > 0 {

		var matched bool
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
				matched = true
				if fm.Merging {
					res.Merge(response)
					break
				}
				res = response
				break
			}
		}

		// break loop if there's no match evaluation in one cycle
		if !matched {
			break
		}
		fm.MaxCycle--
	}

	return

}
