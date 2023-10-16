package rule_time

import (
	"encoding/json"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
)

type After struct {
	Type  string       `json:"type"`
	Left  TimeConstItf `json:"left"`
	Right TimeConstItf `json:"right"`
	hash  *string
}

func (r *After) GetType() string {
	return "rules.time.after"
}

func (r *After) New() rules.RulesItf {
	return new(After)
}

func (r *After) GetHash() string {
	for r.hash == nil {
		hash := common.CreateHash(r.Type, r.Left.GetHash(), r.Right.GetHash())
		r.hash = &hash
	}
	return *r.hash
}

func (r *After) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
	if ok := facts.GetCacheInstance().TryGet(r.GetHash(), &isMatch); ok {
		return isMatch, nil
	}
	defer func() {
		facts.GetCacheInstance().SetCache(r.GetHash(), isMatch)
	}()

	_timeLeft, err := r.Left.GetTime(facts)
	if err != nil {
		return false, nil
	}
	_timeRight, err := r.Right.GetTime(facts)
	if err != nil {
		return false, nil
	}

	return _timeLeft.After(_timeRight), nil
}

func (r *After) UnmarshalJSON(data []byte) (err error) {
	var m map[string]json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}

	for k, val := range m {
		var typeChecker common.TypeChecker

		switch k {
		case "type":
			if err := json.Unmarshal(val, &r.Type); err != nil {
				return err
			}
		case "left":
			var instance TimeConstItf

			if err := json.Unmarshal(val, &typeChecker); err != nil {
				return err
			}
			timeType := Get(typeChecker.Type)
			instance = timeType.New()
			if err := json.Unmarshal(val, instance); err != nil {
				return err
			}
			r.Left = instance
		case "right":
			var instance TimeConstItf

			if err := json.Unmarshal(val, &typeChecker); err != nil {
				return err
			}
			timeType := Get(typeChecker.Type)
			instance = timeType.New()
			if err := json.Unmarshal(val, instance); err != nil {
				return err
			}
			r.Right = instance
		}
	}

	return
}
