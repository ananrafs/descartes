package evaluators

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
)

var (
	evaluatorMap map[string]Factory = make(map[string]Factory)
)

type Factory func() EvaluatorsItf

type EvaluatorsItf interface {
	common.TypeCheckerItf
	Eval(facts.FactsItf) EvalResult
}

type EvalResult struct {
	IsMatch bool
	Error   error
	Result  interface{}
}

func Init(evals ...Factory) {
	for _, evaluator := range evals {
		evaluatorMap[evaluator().GetType()] = evaluator
	}
}

func Get(evaluatorType string) (eval EvaluatorsItf) {
	evaluator, ok := evaluatorMap[evaluatorType]
	if ok {
		return evaluator()
	}
	return
}

func GetCatalog() []EvaluatorsItf {
	res := make([]EvaluatorsItf, 0, len(evaluatorMap))
	for _, factory := range evaluatorMap {
		res = append(res, factory())
	}

	return res
}

func (e *EvalResult) Merge(with EvalResult) {
	e.IsMatch = e.IsMatch || with.IsMatch

	switch e.Result.(type) {
	case int:
		if _target, ok := with.Result.(int); ok {
			e.Result = e.Result.(int) + _target
		}
	case string:
		if _target, ok := with.Result.(string); ok {
			e.Result = e.Result.(string) + _target
		}
	case map[string]interface{}:
		if _target, ok := with.Result.(map[string]interface{}); ok {
			currentMap := e.Result.(map[string]interface{})

			for k, v := range _target {
				currentMap[k] = v
			}
			e.Result = currentMap
		}

	}

	if nil == e.Result && nil != with.Result {
		e.Result = with.Result
	}

	if nil == e.Error && nil != with.Error {
		e.Error = with.Error
	}

}
