package group

import (
	"encoding/json"

	"github.com/ananrafs/descartes/engine/evaluators"
	"github.com/ananrafs/descartes/engine/facts"
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

func (fm *FirstMatch) Eval(facts facts.FactsItf) (res evaluators.EvalResult) {
	for _, eval := range fm.Evaluators {
		res = eval.Eval(facts)
		if res.IsMatch {
			return res
		}
	}
	return
}

func (fm *FirstMatch) UnmarshalJSON(data []byte) (err error) {
	var m map[string]json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}

	for k, val := range m {
		switch k {
		case "type":
			if err := json.Unmarshal(val, &fm.EvalType); err != nil {
				return err
			}
		case "evaluators":
			evalGroup := new(EvaluatorGroup)
			if err := json.Unmarshal(val, &evalGroup); err != nil {
				return err
			}
			fm.Evaluators = *evalGroup
		}
	}
	return
}
