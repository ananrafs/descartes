package core

import (
	"fmt"

	"github.com/ananrafs/descartes/cache"
	"github.com/ananrafs/descartes/engine/actions"
	"github.com/ananrafs/descartes/engine/actions/action"
	action_float "github.com/ananrafs/descartes/engine/actions/action/float"
	action_int "github.com/ananrafs/descartes/engine/actions/action/int"
	action_map "github.com/ananrafs/descartes/engine/actions/action/map"
	actionsgroup "github.com/ananrafs/descartes/engine/actions/group"
	"github.com/ananrafs/descartes/engine/evaluators"
	"github.com/ananrafs/descartes/engine/evaluators/evaluator"
	"github.com/ananrafs/descartes/engine/evaluators/group"
	"github.com/ananrafs/descartes/engine/rules"
	rulesgroup "github.com/ananrafs/descartes/engine/rules/group"
	"github.com/ananrafs/descartes/engine/rules/rule"
	rule_array "github.com/ananrafs/descartes/engine/rules/rule/array"
	rule_bool "github.com/ananrafs/descartes/engine/rules/rule/bool"
	rule_int "github.com/ananrafs/descartes/engine/rules/rule/int"
	rule_string "github.com/ananrafs/descartes/engine/rules/rule/string"
	rule_time "github.com/ananrafs/descartes/engine/rules/rule/time"
	rule_time_type "github.com/ananrafs/descartes/engine/rules/rule/time/time_type"
	"github.com/ananrafs/descartes/law"
)

type Factory struct {
	RuleCreateFunction
	EvalCreateFunction
	ActionCreateFunction
	CacheCreateFunction
	TimeTypeCreateFunction
}

type RuleCreateFunction func() []rules.RulesItf
type EvalCreateFunction func() []evaluators.EvaluatorsItf
type ActionCreateFunction func() []actions.ActionsItf
type CacheCreateFunction func() []cache.CacheItf
type TimeTypeCreateFunction func() []rule_time.TimeConstItf

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
		if factory.TimeTypeCreateFunction != nil {
			rule_time.Init(factory.TimeTypeCreateFunction()...)
		}
	}
}

func Register(laws ...law.Law) error {
	for _, law := range laws {
		_, ok := lawDictionary[law.Slug]
		if ok {
			return fmt.Errorf("law %s already registered", law.Slug)
		}
		lawDictionary[law.Slug] = law
	}

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

func InitTimeType(funcs ...TimeTypeCreateFunction) {
	timeTypes := make([]rule_time.TimeConstItf, 0)
	for _, createFunc := range funcs {
		timeTypes = append(timeTypes, createFunc()...)
	}
	rule_time.Init(timeTypes...)
}

func WithDefaults() Factory {
	return Factory{
		RuleCreateFunction:     WithDefaultRules,
		EvalCreateFunction:     WithDefaultEvaluators,
		CacheCreateFunction:    WithDefaultCaches,
		ActionCreateFunction:   WithDefaultActions,
		TimeTypeCreateFunction: WithDefaultTimeType,
	}
}

func WithDefaultRules() []rules.RulesItf {
	return []rules.RulesItf{
		(&rulesgroup.ConditionalAnd{}).New(),
		(&rulesgroup.ConditionalOr{}).New(),
		(&rulesgroup.ConditionalNot{}).New(),
		(&rule_int.Between{}).New(),
		(&rule_int.Equal{}).New(),
		(&rule_int.Greater{}).New(),
		(&rule_int.Lesser{}).New(),
		(&rule_int.BetweenDynamic{}).New(),
		(&rule_int.EqualDynamic{}).New(),
		(&rule_int.GreaterDynamic{}).New(),
		(&rule_int.LesserDynamic{}).New(),
		(&rule_string.Equal{}).New(),
		(&rule_string.EqualDynamic{}).New(),
		(&rule_string.EqualFold{}).New(),
		(&rule_bool.Bool{}).New(),
		(&rule_array.ArrayContains{}).New(),
		(&rule.Exist{}).New(),
		(&rule.RuleDefault{}).New(),
	}
}

func WithDefaultEvaluators() []evaluators.EvaluatorsItf {
	return []evaluators.EvaluatorsItf{
		(&evaluator.Evaluator{}).New(),
		(&group.FirstMatch{}).New(),
		(&group.MultiMatch{}).New(),
		(&group.MultiMatchOrdered{}).New(),
		(&group.MultiMatchOrderedCycle{}).New(),
		(&evaluator.IterateEvaluator{}).New(),
	}
}

func WithDefaultActions() []actions.ActionsItf {
	return []actions.ActionsItf{
		(&action.Action{}).New(),
		(&actionsgroup.ActionGroup{}).New(),

		(&action_int.Divide{}).New(),
		(&action_int.Mod{}).New(),
		(&action_int.Multiple{}).New(),
		(&action_int.Substract{}).New(),
		(&action_int.Sum{}).New(),

		(&action_float.Divide{}).New(),
		(&action_float.Multiple{}).New(),
		(&action_float.Substract{}).New(),
		(&action_float.Sum{}).New(),

		(&action_map.Append{}).New(),
	}
}

func WithDefaultCaches() []cache.CacheItf {
	return []cache.CacheItf{
		(&cache.Cache{}).New(),
		(&cache.NopCache{}).New(),
	}
}

func WithDefaultTimeType() []rule_time.TimeConstItf {
	return []rule_time.TimeConstItf{
		(&rule_time_type.TimeTypeAddDay{}).New(),
		(&rule_time_type.TimeTypeAddMonth{}).New(),
		(&rule_time_type.TimeTypeAddYear{}).New(),
		(&rule_time_type.TimeTypeDynamic{}).New(),
	}
}
