package core

import (
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

type RuleCreateFunction func() []rules.Factory
type EvalCreateFunction func() []evaluators.Factory
type ActionCreateFunction func() []actions.Factory
type CacheCreateFunction func() []cache.Factory
type TimeTypeCreateFunction func() []rule_time.Factory

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
		lawDictionary[law.Slug] = law
	}

	return nil
}

func InitRule(funcs ...RuleCreateFunction) {
	ruleList := make([]rules.Factory, 0)
	for _, ruleInstance := range funcs {
		ruleList = append(ruleList, ruleInstance()...)
	}
	rules.Init(ruleList...)
}

func InitEvaluator(funcs ...EvalCreateFunction) {
	evalList := make([]evaluators.Factory, 0)
	for _, createFunc := range funcs {
		evalList = append(evalList, createFunc()...)
	}
	evaluators.Init(evalList...)
}

func InitActions(funcs ...ActionCreateFunction) {
	actionList := make([]actions.Factory, 0)
	for _, createFunc := range funcs {
		actionList = append(actionList, createFunc()...)
	}
	actions.Init(actionList...)
}

func InitCaches(funcs ...CacheCreateFunction) {
	caches := make([]cache.Factory, 0)
	for _, createFunc := range funcs {
		caches = append(caches, createFunc()...)
	}
	cache.Init(caches...)
}

func InitTimeType(funcs ...TimeTypeCreateFunction) {
	timeTypes := make([]rule_time.Factory, 0)
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

func WithDefaultRules() []rules.Factory {
	return []rules.Factory{
		rulesgroup.NewConditionalAnd,
		rulesgroup.NewConditionalOr,
		rulesgroup.NewConditionalNot,
		rule_int.NewBetween,
		rule_int.NewEqual,
		rule_int.NewGreater,
		rule_int.NewGreaterEqual,
		rule_int.NewLesser,
		rule_int.NewLesserEqual,
		rule_int.NewBetweenDynamic,
		rule_int.NewEqualDynamic,
		rule_int.NewGreaterDynamic,
		rule_int.NewGreaterEqualDynamic,
		rule_int.NewLesserDynamic,
		rule_int.NewLesserEqualDynamic,
		rule_string.NewEqual,
		rule_string.NewEqualDynamic,
		rule_string.NewEqualFold,
		rule_bool.NewBool,
		rule_array.NewArrayContains,
		rule.NewExist,
		rule.NewRuleDefault,
	}
}

func WithDefaultEvaluators() []evaluators.Factory {
	return []evaluators.Factory{
		evaluator.NewEvaluator,
		group.NewFirstMatch,
		group.NewMultiMatch,
		group.NewMultiMatchOrdered,
		group.NewMultiMatchOrderedCycle,
		evaluator.NewIterateEvaluator,
	}
}

func WithDefaultActions() []actions.Factory {
	return []actions.Factory{
		action.NewAction,
		actionsgroup.NewActionGroup,
		action_int.NewDivide,
		action_int.NewMod,
		action_int.NewMultiple,
		action_int.NewSubtract,
		action_int.NewSum,
		action_float.NewDivide,
		action_float.NewMultiple,
		action_float.NewSubtract,
		action_float.NewSum,
		action_map.NewAppend,
	}
}

func WithDefaultCaches() []cache.Factory {
	return []cache.Factory{
		cache.NewCache,
		cache.NewNopCache,
	}
}

func WithDefaultTimeType() []rule_time.Factory {
	return []rule_time.Factory{
		rule_time_type.NewTimeTypeAddDay,
		rule_time_type.NewTimeTypeAddMonth,
		rule_time_type.NewTimeTypeAddYear,
		rule_time_type.NewTimeTypeDynamic,
	}
}
