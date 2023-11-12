package rule_time_type

import (
	"fmt"
	"time"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	rule_time "github.com/ananrafs/descartes/engine/rules/rule/time"
)

type TimeTypeDynamic struct {
	Type string `json:"string"`
	From string `json:"from"`
	hash *string
}

func (r *TimeTypeDynamic) GetType() string {
	return "time_type.dynamic"
}

func NewTimeTypeDynamic() rule_time.TimeConstItf {
	o := new(TimeTypeDynamic)
	o.Type = o.GetType()
	return o
}

func (r *TimeTypeDynamic) GetHash() string {
	for r.hash == nil {
		hash := common.CreateHash(r.Type, r.From)
		r.hash = &hash
	}
	return *r.hash
}

func (r *TimeTypeDynamic) GetTime(facts facts.FactsItf) (time.Time, error) {

	param := facts.GetMap()
	fields, valStr := new(string), r.From
	if ok := common.ParseFromMustacheTemplate(r.From, fields); ok {
		val, ok := param[*fields]
		if !ok {
			return time.Time{}, common.ErrorNotFoundOnMap(r.From)
		}
		valStr = fmt.Sprintf("%v", val)
	}

	_time, err := time.Parse(rule_time.SUPPORTED_TIME_FORMAT, valStr)
	if err != nil {
		return time.Time{}, err
	}

	return _time, nil
}
