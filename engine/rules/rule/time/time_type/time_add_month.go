package rule_time_type

import (
	"encoding/json"
	"time"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	rule_time "github.com/ananrafs/descartes/engine/rules/rule/time"
)

type TimeTypeAddMonth struct {
	Type   string                 `json:"string"`
	Factor int                    `json:"factor"`
	From   rule_time.TimeConstItf `json:"from"`
	hash   *string
}

func (r *TimeTypeAddMonth) GetType() string {
	return "time_type.add_month"
}

func NewTimeTypeAddMonth() rule_time.TimeConstItf {
	o := new(TimeTypeAddMonth)
	o.Type = o.GetType()
	return o
}

func (r *TimeTypeAddMonth) GetHash() string {
	for r.hash == nil {
		hash := common.CreateHash(r.Type, r.Factor, r.From.GetHash())
		r.hash = &hash
	}
	return *r.hash
}

func (r *TimeTypeAddMonth) GetTime(facts facts.FactsItf) (time.Time, error) {
	_time, err := r.From.GetTime(facts)
	if err != nil {
		return time.Time{}, err
	}
	return _time.AddDate(0, 0, r.Factor), nil
}

func (r *TimeTypeAddMonth) UnmarshalJSON(data []byte) (err error) {
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
		case "factor":
			if err := json.Unmarshal(val, &r.Factor); err != nil {
				return err
			}
		case "from":
			if err := json.Unmarshal(val, &typeChecker); err != nil {
				return err
			}

			instance := rule_time.Get(typeChecker.Type)
			if err := json.Unmarshal(val, instance); err != nil {
				return err
			}
			r.From = instance
		}
	}

	return
}
