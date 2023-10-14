package law

import (
	"encoding/json"
	"fmt"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/evaluators"
)

type Law struct {
	Slug      string                   `json:"slug"`
	Evaluator evaluators.EvaluatorsItf `json:"evaluator"`
}

func (l *Law) Judge(param map[string]interface{}) (interface{}, error) {
	evalRes := l.Evaluator.Eval(param)

	if !evalRes.IsMatch {
		return nil, fmt.Errorf("%v not match on %s law", param, l.Slug)
	}
	return evalRes.Result, evalRes.Error
}

func (l *Law) UnmarshalJSON(data []byte) (err error) {
	var m map[string]json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}

	for k, val := range m {
		switch k {
		case "slug":
			if err = json.Unmarshal(val, &l.Slug); err != nil {
				return
			}
		case "evaluator":
			var typeChecker common.TypeChecker
			if err = json.Unmarshal(val, &typeChecker); err != nil {
				return
			}

			eval := evaluators.Get(typeChecker.Type)
			if err = json.Unmarshal(val, eval); err != nil {
				return
			}
			l.Evaluator = eval
		}
	}
	return
}

func CreateLaw(jsonStr string) (l Law, err error) {
	err = json.Unmarshal([]byte(jsonStr), &l)
	if err != nil {
		return
	}
	return
}
