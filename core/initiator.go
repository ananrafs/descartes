package core

import (
	"fmt"

	"github.com/ananrafs/descartes/evaluators"
	"github.com/ananrafs/descartes/evaluators/evaluator"
	"github.com/ananrafs/descartes/evaluators/group"
	"github.com/ananrafs/descartes/law"
	"github.com/ananrafs/descartes/rules"
	rulesgroup "github.com/ananrafs/descartes/rules/group"
	rule_int "github.com/ananrafs/descartes/rules/rule/int"
	rule_string "github.com/ananrafs/descartes/rules/rule/string"
)

type RuleCreateFunction func() []rules.RulesItf
type EvalCreateFunction func() []evaluators.EvaluatorItf

func Register(law law.Law) error {
	_, ok := lawDictionary[law.Slug]
	if ok {
		return fmt.Errorf("law %s already registered", law.Slug)
	}
	lawDictionary[law.Slug] = law

	return nil
}

func InitRule(funcs ...RuleCreateFunction) {
	ruleList := make([]rules.RulesItf, 0)
	for _, ruleInstance := range funcs {
		ruleList = append(ruleList, ruleInstance()...)
	}
	rules.InitRules(ruleList...)
}

func InitEvaluator(funcs ...EvalCreateFunction) {
	evalList := make([]evaluators.EvaluatorItf, 0)
	for _, evalInstance := range funcs {
		evalList = append(evalList, evalInstance()...)
	}
	evaluators.InitEvaluators(evalList...)
}

func WithDefaultRules() []rules.RulesItf {
	return []rules.RulesItf{
		&rulesgroup.ConditionalAnd{},
		&rulesgroup.ConditionalOr{},
		&rule_int.RuleIntBetween{},
		&rule_int.RuleIntEqual{},
		&rule_int.RuleIntGreater{},
		&rule_int.RuleIntLesser{},
		&rule_string.RuleStringEqual{},
		&rule_string.RuleStringEqualFold{},
	}
}

func WithDefaultEvaluators() []evaluators.EvaluatorItf {
	return []evaluators.EvaluatorItf{
		&evaluator.Evaluator{},
		&group.FirstMatch{},
	}
}
