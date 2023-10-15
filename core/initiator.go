package core

import (
	"fmt"

	"github.com/ananrafs/descartes/cache"
	"github.com/ananrafs/descartes/engine/actions"
	"github.com/ananrafs/descartes/engine/actions/action"
	action_int "github.com/ananrafs/descartes/engine/actions/action/int"
	actionsgroup "github.com/ananrafs/descartes/engine/actions/group"
	"github.com/ananrafs/descartes/engine/evaluators"
	"github.com/ananrafs/descartes/engine/evaluators/evaluator"
	"github.com/ananrafs/descartes/engine/evaluators/group"
	"github.com/ananrafs/descartes/engine/rules"
	rulesgroup "github.com/ananrafs/descartes/engine/rules/group"
	"github.com/ananrafs/descartes/engine/rules/rule"
	rule_bool "github.com/ananrafs/descartes/engine/rules/rule/bool"
	rule_int "github.com/ananrafs/descartes/engine/rules/rule/int"
	rule_string "github.com/ananrafs/descartes/engine/rules/rule/string"
	"github.com/ananrafs/descartes/law"
)

type Factory struct {
	RuleCreateFunction
	EvalCreateFunction
	ActionCreateFunction
	CacheCreateFunction
}

type RuleCreateFunction func() []rules.RulesItf
type EvalCreateFunction func() []evaluators.EvaluatorsItf
type ActionCreateFunction func() []actions.ActionsItf
type CacheCreateFunction func() []cache.CacheItf

func InitFactory(factories ...Factory) {
	for _, factory := range factories {
		if factory.RuleCreateFunction != nil {
			rules.Init(factory.RuleCreateFunction()...)
		}
		if factory.EvalCreateFunction != nil {
			evaluators.Init(factory.EvalCreateFunction()...)
		}
		if factory.ActionCreateFunction != nil {
			actions.Init(factory.ActionCreateFunction()...)
		}
		if factory.CacheCreateFunction != nil {
			cache.Init(factory.CacheCreateFunction()...)
		}
	}
}

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

func InitCaches(funcs ...CacheCreateFunction) {
	caches := make([]cache.CacheItf, 0)
	for _, createFunc := range funcs {
		caches = append(caches, createFunc()...)
	}
	cache.Init(caches...)
}

func WithDefaults() Factory {
	return Factory{
		RuleCreateFunction:   WithDefaultRules,
		EvalCreateFunction:   WithDefaultEvaluators,
		CacheCreateFunction:  WithDefaultCaches,
		ActionCreateFunction: WithDefaultActions,
	}
}

func WithDefaultRules() []rules.RulesItf {
	return []rules.RulesItf{
		&rulesgroup.ConditionalAnd{},
		&rulesgroup.ConditionalOr{},
		&rulesgroup.ConditionalNot{},
		&rule_int.RuleIntBetween{},
		&rule_int.RuleIntEqual{},
		&rule_int.RuleIntGreater{},
		&rule_int.RuleIntLesser{},
		&rule_int.RuleIntBetweenDynamic{},
		&rule_int.RuleIntEqualDynamic{},
		&rule_int.RuleIntGreaterDynamic{},
		&rule_int.RuleIntLesserDynamic{},
		&rule_string.RuleStringEqual{},
		&rule_string.RuleStringEqualFold{},
		&rule_bool.RuleBool{},
		&rule.RuleDefault{},
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
		&actionsgroup.ActionGroup{},
		&action_int.ActionIntDivide{},
		&action_int.ActionIntMod{},
		&action_int.ActionIntMultiple{},
		&action_int.ActionIntSubstract{},
		&action_int.ActionIntSum{},
	}
}

func WithDefaultCaches() []cache.CacheItf {
	return []cache.CacheItf{
		&cache.Cache{},
		&cache.NopCache{},
	}
}
