package law

import (
	"encoding/json"

	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
)

type Fact struct {
	Slug  string         `json:"slug"`
	Facts facts.FactsItf `json:"param"`
}

func (f *Fact) UnmarshalJSON(data []byte) (err error) {
	var m map[string]json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}

	for k, val := range m {
		switch k {
		case "slug":
			if err = json.Unmarshal(val, &f.Slug); err != nil {
				return
			}
		case "param":
			facts := new(facts.Facts)
			if err = json.Unmarshal(data, facts); err != nil {
				return
			}

			f.Facts = facts
		}
	}
	return
}

func (b *Fact) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Slug  string                 `json:"slug"`
		Param map[string]interface{} `json:"param"`
	}{
		Slug:  b.Slug,
		Param: b.Facts.GetMap(),
	})
}

func CreateFact(jsonStr string) (f Fact, err error) {
	err = json.Unmarshal([]byte(jsonStr), &f)
	if err != nil {
		return
	}

	return f, nil
}

func CreateMultipleFact(jsonStr string) (fs []Fact, err error) {
	err = json.Unmarshal([]byte(jsonStr), &fs)
	if err != nil {
		return
	}

	return fs, nil
}

// FactMaker
//
//	to use :
//	fact := MakeFact(yourmap)
//	.AddFields(anothermap)
//	...
//	.Generate(yourlovelyslug)
type FactMaker func(*map[string]interface{})

func MakeFact(mps ...map[string]interface{}) FactMaker {
	return func(m *map[string]interface{}) {
		for _, mp := range mps {
			*m = common.ManipulateMap(*m).Merge(mp)
		}
	}
}

// merge given map to current map
func (fmake FactMaker) AddFields(mp map[string]interface{}) FactMaker {
	return func(m *map[string]interface{}) {
		fmake(m)
		*m = common.ManipulateMap(*m).Merge(mp)
	}
}

func (fmake FactMaker) Generate(slug string) Fact {
	var (
		_default = make(map[string]interface{})
	)
	fmake(&_default)

	return Fact{
		Slug: slug,
		Facts: &facts.Facts{
			Fields: common.ManipulateMap().DeepCopy(_default),
		},
	}
}
