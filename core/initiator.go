package core

import (
	"fmt"

	"github.com/ananrafs/descartes/engine/actions"
	"github.com/ananrafs/descartes/engine/actions/action"
	action_int "github.com/ananrafs/descartes/engine/actions/action/int"
	"github.com/ananrafs/descartes/engine/evaluators"
	"github.com/ananrafs/descartes/engine/evaluators/evaluator"
	"github.com/ananrafs/descartes/engine/evaluators/group"
	"github.com/ananrafs/descartes/engine/rules"
	rulesgroup "github.com/ananrafs/descartes/engine/rules/group"
	rule_int "github.com/ananrafs/descartes/engine/rules/rule/int"
	rule_string "github.com/ananrafs/descartes/engine/rules/rule/string"
	"github.com/ananrafs/descartes/law"
)

type RuleCreateFunction func() []rules.RulesItf
type EvalCreateFunction func() []evaluators.EvaluatorsItf
type ActionCreateFunction func() []actions.ActionsItf

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
	rules.Init(ruleList...)
}

func InitEvaluator(funcs ...EvalCreateFunction) {
	evalList := make([]evaluators.EvaluatorsItf, 0)
	for _, createFunc := range funcs {
		evalList = append(evalList, createFunc()...)
	}
	evaluators.Init(evalList...)
}

func InitActions(funcs ...ActionCreateFunction) {
	actionList := make([]actions.ActionsItf, 0)
	for _, createFunc := range funcs {
		actionList = append(actionList, createFunc()...)
	}
	actions.Init(actionList...)
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

func WithDefaultEvaluators() []evaluators.EvaluatorsItf {
	return []evaluators.EvaluatorsItf{
		&evaluator.Evaluator{},
		&group.FirstMatch{},
	}
}

func WithDefaultActions() []actions.ActionsItf {
	return []actions.ActionsItf{
		&action.Action{},
		&action_int.ActionIntDivide{},
		&action_int.ActionIntMod{},
		&action_int.ActionIntMultiple{},
		&action_int.ActionIntSubstract{},
		&action_int.ActionIntSum{},
	}
}
