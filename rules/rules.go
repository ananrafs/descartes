package rules

import "strings"

var (
	ruleList []RulesItf
)

type RulesItf interface {
	New() RulesItf
	GetRules() string
	IsMatch(map[string]interface{}) bool
}

func GetRules(rulesType string) (rule RulesItf) {
	for _, rulee := range ruleList {
		if strings.EqualFold(rulee.GetRules(), rulesType) {
			rule = rulee
		}
	}
	return
}

func InitRules(rules ...RulesItf) {
	ruleList = append(ruleList, rules...)
}
