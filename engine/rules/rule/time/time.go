package rule_time

import (
	"time"

	"github.com/ananrafs/descartes/engine/facts"
)

const (
	SUPPORTED_TIME_FORMAT = "2006-01-02 15:04:05"
)

type TimeConstItf interface {
	GetType() string
	New() TimeConstItf
	GetHash() string

	GetTime(facts.FactsItf) (time.Time, error)
}

var (
	timeTypeMap map[string]TimeConstItf = make(map[string]TimeConstItf)
)

func Init(timeTypes ...TimeConstItf) {
	for _, timeType := range timeTypes {
		timeTypeMap[timeType.GetType()] = timeType
	}
}

func Get(timeTypeKey string) (timeType TimeConstItf) {
	timeType, ok := timeTypeMap[timeTypeKey]
	if ok {
		return timeType
	}

	return
}
