package evaluator

import (
	"encoding/json"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/actions"
	"github.com/ananrafs/descartes/engine/evaluators"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type Evaluator struct {
	EvaluatorType string             `json:"type"`
	Rules         rules.RulesItf     `json:"rule"`
	Action        actions.ActionsItf `json:"action"`
}

func (e *Evaluator) GetType() string {
	return "evaluator"
}

func (e *Evaluator) New() evaluators.EvaluatorsItf {
	return new(Evaluator)
}

func (e *Evaluator) Eval(fact facts.FactsItf) (res evaluators.EvalResult) {
	isMatch, err := e.Rules.IsMatch(fact)
	if !isMatch {
		res.Error = err
		return
	}
	res.IsMatch = isMatch
	actionRes, err := e.Action.Do(fact)
	if err != nil {
		res.Error = err
		return
	}
	res.Result = actionRes

	return
}

func (e *Evaluator) UnmarshalJSON(data []byte) (err error) {
	var m map[string]json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}

	for k, val := range m {
		var typeChecker common.TypeChecker

		switch k {
		case "type":
			if err := json.Unmarshal(val, &e.EvaluatorType); err != nil {
				return err
			}
		case "action":
			var instance actions.ActionsItf

			if err := json.Unmarshal(val, &typeChecker); err != nil {
				return err
			}

			action := actions.Get(typeChecker.Type)
			instance = action.New()
			if err := json.Unmarshal(val, instance); err != nil {
				return err
			}
			e.Action = instance
		case "rule":
			var instance rules.RulesItf

			if err := json.Unmarshal(val, &typeChecker); err != nil {
				return err
			}
			rule := rules.Get(typeChecker.Type)
			instance = rule.New()
			if err := json.Unmarshal(val, instance); err != nil {
				return err
			}

			e.Rules = instance
		}
	}

	return
}
