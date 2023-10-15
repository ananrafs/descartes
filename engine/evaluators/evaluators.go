package evaluators

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
)

var (
	evaluatorMap map[string]EvaluatorsItf = make(map[string]EvaluatorsItf)
)

type EvaluatorsItf interface {
	common.TypeCheckerItf
	New() EvaluatorsItf

	Eval(facts.FactsItf) EvalResult
}

type EvalResult struct {
	IsMatch bool
	Error   error
	Result  interface{}
}

func Init(evals ...EvaluatorsItf) {
	for _, evaluator := range evals {
		evaluatorMap[evaluator.GetType()] = evaluator
	}
}

func Get(evaluatorType string) (evaluator EvaluatorsItf) {
	evaluator, ok := evaluatorMap[evaluatorType]
	if ok {
		return evaluator
	}
	return
}
