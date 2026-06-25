package rule_array

import (
	"github.com/ananrafs/descartes/common"
	"github.com/ananrafs/descartes/engine/facts"
	f "github.com/ananrafs/descartes/engine/facts"
	"github.com/ananrafs/descartes/engine/rules"
	json "github.com/json-iterator/go"
)

type StructContains struct {
	Type  string         `json:"type"`
	Field string         `json:"field"`
	Rule  rules.RulesItf `json:"rule"`

	hash *string
}

func (s *StructContains) GetType() string {
	return "rules.array.struct.contains"
}

func NewStructContains() rules.RulesItf {
	o := new(StructContains)
	o.Type = o.GetType()
	return o
}

func (s *StructContains) GetHash() string {
	for s.hash == nil {
		hash := common.CreateHash(s.Type, s.Field, s.Rule)
		s.hash = &hash
	}
	return *s.hash
}

func (s *StructContains) IsMatch(facts facts.FactsItf) (isMatch bool, err error) {
	if ok := facts.GetCacheInstance().TryGet(s.GetHash(), &isMatch); ok {
		return isMatch, nil
	}
	defer func() {
		facts.GetCacheInstance().SetCache(s.GetHash(), isMatch)
	}()
	param := facts.GetMap()
	v, ok := param[s.Field]
	if !ok {
		return false, common.ErrorNotFoundOnMap(s.Field)
	}

	intv := new([]interface{})
	if err = common.ConvertToArray(v, intv); err != nil {
		return false, err
	}

	for _, val := range *intv {
		factField, ok := val.(map[string]interface{})
		if !ok {
			return false, common.ErrorCasting(val)
		}
		obj := &f.Facts{
			Fields: common.ManipulateMap().Copy(factField),
		}

		obj.SetCacheInstance(facts.GetCacheInstance())
		match, err := s.Rule.IsMatch(obj)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}

	return false, nil
}

func (s *StructContains) UnmarshalJSON(data []byte) (err error) {
	var m map[string]json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}

	for k, val := range m {
		var typeChecker common.TypeChecker

		switch k {
		case "type":
			if err := json.Unmarshal(val, &s.Type); err != nil {
				return err
			}
		case "field":
			if err := json.Unmarshal(val, &s.Field); err != nil {
				return err
			}
		case "rule":
			var instance rules.RulesItf
			if err := json.Unmarshal(val, &typeChecker); err != nil {
				return err
			}
			instance = rules.Get(typeChecker.Type)
			if err := json.Unmarshal(val, instance); err != nil {
				return err
			}
			s.Rule = instance
		}
	}

	return
}
