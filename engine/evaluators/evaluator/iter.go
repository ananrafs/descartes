package evaluator

import (
	"encoding/json"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/evaluators"
	"github.com/ananrafs/descartes/engine/facts"
)

type IterateEvaluator struct {
	EvaluatorType string                   `json:"type"`
	Field         string                   `json:"field"`
	Iterant       string                   `json:"iterant"`
	Evaluator     evaluators.EvaluatorsItf `json:"evaluator"`
}

func (e *IterateEvaluator) GetType() string {
	return "evaluator.iterate"
}

func NewIterateEvaluator() evaluators.EvaluatorsItf {
	o := new(IterateEvaluator)
	o.EvaluatorType = o.GetType()
	return o
}

func (e IterateEvaluator) Eval(fact facts.FactsItf) (res evaluators.EvalResult) {
	param := fact.GetMap()

	mapVal, ok := param[e.Iterant]
	if !ok {
		return
	}

	iteratedObjects, ok := mapVal.([]interface{})
	if !ok {
		return
	}

	for _, iteratedObject := range iteratedObjects {
		param[e.Field] = iteratedObject
		response := e.Evaluator.Eval(fact)
		res.Merge(response)
	}

	return
}

func (e *IterateEvaluator) UnmarshalJSON(data []byte) (err error) {
	var m map[string]json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}

	for k, val := range m {
		switch k {
		case "type":
			if err := json.Unmarshal(val, &e.EvaluatorType); err != nil {
				return err
			}
		case "field":
			if err = json.Unmarshal(val, &e.Field); err != nil {
				return
			}
		case "iterant":
			if err = json.Unmarshal(val, &e.Iterant); err != nil {
				return
			}
		case "evaluator":
			var typeChecker common.TypeChecker
			if err = json.Unmarshal(val, &typeChecker); err != nil {
				return
			}

			instance := evaluators.Get(typeChecker.Type)
			if err = json.Unmarshal(val, instance); err != nil {
				return
			}

			e.Evaluator = instance
		}
	}

	return
}
