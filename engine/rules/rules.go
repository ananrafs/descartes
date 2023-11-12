package rules

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
)

var (
	ruleMap map[string]Factory = make(map[string]Factory)
)

type Factory func() RulesItf

type RulesItf interface {
	GetHash() string
	common.TypeCheckerItf

	IsMatch(facts.FactsItf) (bool, error)
}

func Init(rules ...Factory) {
	for _, rule := range rules {
		ruleMap[rule().GetType()] = rule
	}
}

func Get(rulesType string) (rule RulesItf) {
	ruleFactory, ok := ruleMap[rulesType]
	if ok {
		return ruleFactory()
	}

	return
}

func GetCatalog() []RulesItf {
	res := make([]RulesItf, 0, len(ruleMap))
	for _, factory := range ruleMap {
		res = append(res, factory())
	}

	return res
}
