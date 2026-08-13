package bench

import (
	"github.com/ananrafs/descartes/core"
	"github.com/ananrafs/descartes/law"
	"fmt"
	"log"
	"testing"
)

var jsonRule = `{"slug":"section_layout_platform_fee_plus","evaluator":{"type":"evaluator.group.first_match","evaluators":[{"type":"evaluator","name":"PLATFORM_FEE_PLUS","rule":{"type":"rules.conditional.and","rules":[{"type":"rules.array.contains","value":"PLATFORM_FEE_PLUS","field":"user_benefits"}]},"action":{"is_eligible":"true"}},{"type":"evaluator","name":"default","rule":{"type":"rules.default"},"action":{"is_eligible":"false"}}]}}`
var prettyJsonRule = `
{
    "slug": "section_layout_platform_fee_plus",
    "evaluator":
    {
        "type": "evaluator.group.first_match",
        "evaluators":
        [
            {
                "type": "evaluator",
                "name": "PLATFORM_FEE_PLUS",
                "rule":
                {
                    "type": "rules.conditional.and",
                    "rules":
                    [
                        {
                            "type": "rules.array.contains",
                            "value": "PLATFORM_FEE_PLUS",
                            "field": "user_benefits"
                        }
                    ]
                },
                "action":
                {
                    "is_eligible": "true"
                }
            },
            {
                "type": "evaluator",
                "name": "default",
                "rule":
                {
                    "type": "rules.default"
                },
                "action":
                {
                    "is_eligible": "false"
                }
            }
        ]
    }
}
`

var fact = map[string]interface{}{
	"user_segments": []interface{}{"some_segment", "some2_segment"},
	"user_benefits": []interface{}{"BBO_PLUS", "PLATFORM_FEE_PLUS", "GF_PLUS"},
}

var factWithArrayString = map[string]interface{}{
	"user_segments": []string{"some_segment", "some2_segment"},
	"user_benefits": []string{"BBO_PLUS", "PLATFORM_FEE_PLUS", "GF_PLUS"},
}

func SimulateRule(lawString string, fact map[string]interface{}) (result map[string]interface{}, err error) {
	tempLaw, err := law.CreateLaw(lawString)
	if err != nil {
		return result, err
	}

	tempFact := law.MakeFact(fact).Generate(tempLaw.Slug)

	res, err := tempLaw.Judge(tempFact.Facts)
	if err != nil {
		return result, err
	}
	resMap, ok := res.(map[string]interface{})
	if !ok {
		err = fmt.Errorf("Unrecognized result")
		return result, err
	}

	return resMap, nil
}

func BenchmarkJudge(b *testing.B) {
	core.InitFactory(core.WithDefaults())
	b.Run("EvalRegister", func(b *testing.B) {
		tempLaw, _ := law.CreateLaw(jsonRule)
		_ = core.Register(tempLaw)
		for i := 0; i < b.N; i++ {
			tempFact := law.MakeFact(fact).Generate(tempLaw.Slug)
			_, err := core.Eval(tempFact)
			if err != nil {
				log.Fatalf("err: %s\n", err.Error())
			}
		}
	})

	b.Run("EvalSimulate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := SimulateRule(jsonRule, fact)

			if err != nil {
				log.Fatalf("err: %s\n", err.Error())
			}
		}
	})

	b.Run("EvalSimulateWithPrettyJson", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := SimulateRule(prettyJsonRule, fact)

			if err != nil {
				log.Fatalf("err: %s\n", err.Error())
			}
		}
	})

	b.Run("EvalRegister_FactWithArrayString", func(b *testing.B) {
		tempLaw, _ := law.CreateLaw(jsonRule)
		_ = core.Register(tempLaw)
		for i := 0; i < b.N; i++ {
			tempFact := law.MakeFact(factWithArrayString).Generate(tempLaw.Slug)
			_, err := core.Eval(tempFact)
			if err != nil {
				log.Fatalf("err: %s\n", err.Error())
			}
		}
	})

	b.Run("EvalSimulate_FactWithArrayString", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := SimulateRule(jsonRule, factWithArrayString)

			if err != nil {
				log.Fatalf("err: %s\n", err.Error())
			}
		}
	})

	b.Run("EvalSimulateWithPrettyJson_FactWithArrayString", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := SimulateRule(prettyJsonRule, factWithArrayString)

			if err != nil {
				log.Fatalf("err: %s\n", err.Error())
			}
		}
	})
}
