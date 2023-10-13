package evaluators

import (
	"strings"
)

var (
	evaluatorList []EvaluatorItf
)

func GetEvaluators(evaluatorType string) (evaluator EvaluatorItf) {
	for _, eval := range evaluatorList {
		if strings.EqualFold(eval.GetType(), evaluatorType) {
			evaluator = eval
		}
	}
	return
}

type EvaluatorItf interface {
	GetType() string
	New() EvaluatorItf
	Eval(param map[string]interface{}) EvalResult
}

type EvalResult struct {
	IsMatch bool
	Error   error
	Result  interface{}
}

func InitEvaluators(evals ...EvaluatorItf) {
	evaluatorList = append(evaluatorList, evals...)
}
