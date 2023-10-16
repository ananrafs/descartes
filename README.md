# Descartes

Descartes is tool to evaluate fact from rules-action ruleset (law).

example `law`:
```
{
    "slug": "geometry_dynamic",
    "evaluator": {
        "type": "evaluator.group.first_match",
        "evaluators": [
            {
                "type": "evaluator",
                "rule": {
                    "type": "rules.int.greater.dynamic",
                    "left": "{{ radius }}",
                    "right": 0
                },
                "action": {
                    "type": "actions.group",
                    "actions": [
                        {
                            "type": "actions.int.multiple",
                            "field": "area",
                            "factors": [
                                2,
                                3,
                                "{{ radius }}"
                            ]
                        },
                        {
                            "area": "{{ area }}",
                            "object": "{{ name }}"
                        }
                    ]
                }
            },
            {
                "type": "evaluator",
                "rule": {
                    "type": "rules.conditional.and",
                    "rules": [
                        {
                            "type": "rules.int.greater.dynamic",
                            "left": "{{ length }}",
                            "right": 0
                        },
                        {
                            "type": "rules.int.greater",
                            "field": "width",
                            "value": 0
                        }
                    ]
                },
                "action": {
                    "type": "actions.group",
                    "actions": [
                        {
                            "type": "actions.int.multiple",
                            "field": "area",
                            "factors": [
                                "{{ length }}",
                                "{{ width }}"
                            ]
                        },
                        {
                            "area": "{{ area }}",
                            "object": "{{ name }}"
                        }
                    ]
                }
            },
            {
                "type": "evaluator",
                "rule": {
                    "type": "rules.int.greater",
                    "field": "side",
                    "value": 0
                },
                "action": {
                    "type": "actions.group",
                    "actions": [
                        {
                            "type": "actions.int.multiple",
                            "field": "area",
                            "factors": [
                                "{{ side }}",
                                "{{ side }}"
                            ]
                        },
                        {
                            "area": "{{ area }}",
                            "object": "{{ name }}"
                        }
                    ]
                }
            },
            {
                "type": "evaluator",
                "rule": {
                    "type": "rules.default"
                },
                "action": {
                    "object": "{{ name }}",
                    "area": "cant calculate"
                }
            }
        ]
    }
}
```

example `facts`
```
[
    {
        "slug": "geometry_dynamic",
        "param": {
            "name": "square",
            "side": 11
        }
    },
    {
        "slug": "geometry_dynamic",
        "param": {
            "name": "rectangle",
            "length": 11,
            "width": 5
        }
    },
    {
        "slug": "geometry_dynamic",
        "param": {
            "name": "circle",
            "radius": 12
        }
    },
    {
        "slug": "geometry_dynamic",
        "param": {
            "name": "gundam",
            "grade": "master grade"
        }
    }
]
```

example `output` :
```
[
	{
		"area": 121,
		"object": "square"
	},
	{
		"area": 55,
		"object": "rectangle"
	},
	{
		"area": 72,
		"object": "circle"
	},
	{
		"area": "cant calculate",
		"object": "gundam"
	}
]
```
## 

## Getting Started
- init 
```
// init default types
core.InitFactory(core.WithDefaults())

```
- create law, and register it
```
l, err := law.CreateLaw(jsonLaw)
if err != nil {
	panic(err)
}

// register
err = core.Register(l)
if err != nil {
	panic(err)
}
```
- make a fact
```
fact, err := law.CreateFact(jsonFact)
if err != nil {
	panic(err)
}
```

- evaluate fact based on registered law
```
res, err := core.Eval(fact)
if err != nil {
	panic(err)
}		
```

full code :

	package main

	import (
		"github.com/ananrafs/descartes/core"
		"github.com/ananrafs/descartes/law"
	)


	func main() {
		// init default types
		core.InitFactory(core.WithDefaults())

		var (
			jsonLaw string
			jsonFact string
		)

		// put your json law
		l, err := law.CreateLaw(jsonLaw)
		if err != nil {
			panic(err)
		}

		// register
		err = core.Register(l)
		if err != nil {
			panic(err)
		}

		fact, err := law.CreateFact(jsonFact)
		if err != nil {
			panic(err)
		}

		//evaluate fact based on registered law
		res, err := core.Eval(fact)
		if err != nil {
			panic(err)
		}		
	}
	
	

## Structure

### Law
Law structured with 2 fields, slug and evaluator.
```
{
    "slug": "{{ slug }}",
    "evaluator": {
       "{{ evaluator }}"
    }
}
```
each law contains of `slug` as identifier, and `evaluator` as a law-evaluator

### Evaluator
Evaluator structured with at least 1 fields, which is `type` as identifier
```
 {
        "type": " {{ type }}",
		...
 }
```
Evaluator working as a translator of `Fact` to `Result`. Currently there are 2 Evaluators, `evaluator` and `evaluator.group.first_match`, and you can extend it on based on your usecase.

to create new Evaluator :
```
type CustomEvaluator struct {
	// mandatory
	EvaluatorType string             `json:"type"`

	// optional
	Rules         rules.RulesItf     `json:"rule"`
	Action        actions.ActionsItf `json:"action"`
	//... and so on
}

func (e *CustomEvaluator) GetType() string {
	return "evaluator.custom"
}

func (e *CustomEvaluator) New() evaluators.EvaluatorsItf {
	return new(CustomEvaluator)
}

func (e *CustomEvaluator) Eval(fact facts.FactsItf) (res evaluators.EvalResult) {
	// evaluate based on your usecase
	return
}
```

then register it
```
core.InitFactory(
	// defaults and optional, you may use it or create another 
	core.WithDefaults(), 
	core.Factory{
		EvalCreateFunction: func() []evaluators.EvaluatorsItf {
					return []evaluators.EvaluatorsItf{
						// add evaluator
						&CustomEvaluator{},
						// ... and so on
						}
					},
				},
			)
```
and you can make law like this 
```
{
    "slug": "custom-law",
    "evaluator": {
       "type": "evaluator.custom",
	   ...
    }
}
```

### Rules
Rules is used on default `evaluator` to validate and then perform an `action`. same as `Evaluator`, you have to put `type` on your struct

```
 {
    "type": "{{ type }}",
    ...
}

```
to create custom rule:

```
type CustomRule struct {
	// mandatory
	RuleType string `json:"type"`

	// optional
	Field string `json:"field"`

	// private fields used for hashing for cache
	hash     *string
}

func (c *CustomRule) GetType() string {
	return "rules.custom"
}

func (c *CustomRule) New() rules.RulesItf {
	return new(CustomRule)
}

func (c *CustomRule) GetHash() string {
	for c.hash == nil {
		hash := common.CreateHash(c.RuleType,...)
		c.hash = &hash
	}
	return *c.hash
}


func (c *RuleDefault) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
	return c.Field == "ruok", nil
}

```
and then register it
```
core.InitFactory(
	// defaults and optional, you may use it or create another 
	core.WithDefaults(), 
	core.Factory{
		RuleCreateFunction: func() []rules.RulesItf {
			return []rules.RulesItf{
				&CustomRule{},
				// ... and so on
			}
		},
	},
)
```
and you can use your `rule.custom` like this :
```
{
    "slug": "custom-law",
    "evaluator": {
       "type": "evaluator.custom",
	   "rule": {
			"type": "rule.custom",
			"field":"ruok"
	   }
    }
}
```