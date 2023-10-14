package law

import (
	"encoding/json"
)

type Fact struct {
	Slug  string                 `json:"slug"`
	Param map[string]interface{} `json:"param"`
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
