package evaluator

import (
	"encoding/json"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/evaluators"
	"github.com/ananrafs/descartes/rules"
)

type Evaluator struct {
	EvaluatorType string         `json:"type"`
	Rules         rules.RulesItf `json:"rule"`
	Action        interface{}    `json:"action"`
}

func (e *Evaluator) GetType() string {
	return "evaluator"
}

func (fm *Evaluator) New() evaluators.EvaluatorItf {
	return new(Evaluator)
}

func (e *Evaluator) Eval(param map[string]interface{}) (res evaluators.EvalResult) {
	isMatch := e.Rules.IsMatch(param)
	res.IsMatch = isMatch
	if res.IsMatch {
		res.Result = e.Action
	}
	return
}

func (r *Evaluator) UnmarshalJSON(data []byte) (err error) {
	var m map[string]json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}

	for k, val := range m {
		switch k {
		case "type":
			if err := json.Unmarshal(val, &r.EvaluatorType); err != nil {
				return err
			}
		case "action":
			if err := json.Unmarshal(val, &r.Action); err != nil {
				return err
			}
		case "rule":
			var typeChecker common.TypeChecker
			if err := json.Unmarshal(val, &typeChecker); err != nil {
				return err
			}

			rule := rules.GetRules(typeChecker.Type)
			newInstance := rule.New()
			if err := json.Unmarshal(val, newInstance); err != nil {
				return err
			}
			r.Rules = newInstance
		}
	}

	return
}
