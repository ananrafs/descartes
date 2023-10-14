package rules

import "github.com/ananrafs/descartes/common"

var (
	ruleMap map[string]RulesItf = make(map[string]RulesItf)
)

type RulesItf interface {
	New() RulesItf
	common.TypeCheckerItf
	IsMatch(map[string]interface{}) (bool, error)
}

func Init(rules ...RulesItf) {
	for _, rule := range rules {
		ruleMap[rule.GetType()] = rule
	}
}

func Get(rulesType string) (rule RulesItf) {
	rule, ok := ruleMap[rulesType]
	if ok {
		return rule
	}

	return
}
