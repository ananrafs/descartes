package rule_time

import (
	"time"

	"github.com/ananrafs/descartes/engine/facts"
)

const (
	SUPPORTED_TIME_FORMAT = "2006-01-02 15:04:05"
)

type Factory func() TimeConstItf

type TimeConstItf interface {
	GetType() string
	GetHash() string

	GetTime(facts.FactsItf) (time.Time, error)
}

var (
	timeTypeMap map[string]Factory = make(map[string]Factory)
)

func Init(timeTypes ...Factory) {
	for _, timeType := range timeTypes {
		timeTypeMap[timeType().GetType()] = timeType
	}
}

func Get(timeTypeKey string) (timeType TimeConstItf) {
	factory, ok := timeTypeMap[timeTypeKey]
	if ok {
		return factory()
	}

	return
}
