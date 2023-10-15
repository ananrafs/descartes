package law

import (
	"encoding/json"

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
