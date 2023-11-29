package law

import (
	"encoding/json"

	"github.com/ananrafs/descartes/cache"
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/evaluators"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/errors"
)

type Law struct {
	Slug      string                   `json:"slug"`
	Evaluator evaluators.EvaluatorsItf `json:"evaluator"`
	Cache     string                   `json:"cache"`
}

func (l Law) Judge(facts facts.FactsItf) (interface{}, error) {
	facts.SetCacheInstance(cache.Get(l.Cache))
	evalRes := l.Evaluator.Eval(facts)

	if !evalRes.IsMatch {
		return nil, errors.ErrFactsNotMatch(l.Slug)
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

			instance := evaluators.Get(typeChecker.Type)
			if err = json.Unmarshal(val, instance); err != nil {
				return
			}

			l.Evaluator = instance
		case "cache":
			if err = json.Unmarshal(val, &l.Cache); err != nil {
				return
			}
		}
	}

	return
}

func CreateLaw(jsonStr string) (l Law, err error) {
	return CreateLawFromJsonByte([]byte(jsonStr))
}

func CreateLawFromJsonByte(jsonByte []byte) (l Law, err error) {
	err = json.Unmarshal(jsonByte, &l)
	if err != nil {
		return
	}
	return
}