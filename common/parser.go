package common

import "strings"

func ParseFromTemplate(source interface{}, dest *string) (isMatch bool) {
	strSource, ok := source.(string)
	if !ok {
		return false
	}

	return GetTemplatedString(strSource, dest)
}

func GetTemplatedString(source string, dest *string) bool {
	source = strings.TrimSpace(source)
	open := "{{"
	close := "}}"

	if source[:2] == open && source[len(source)-2:] == close {
		trimmed := strings.TrimSpace(source[2 : len(source)-2])
		*dest = trimmed
		return true
	}
	return false

}
func DeepParseFromTemplate(source interface{}, dest *string) (isMatch bool, depth int) {
	strSource, ok := source.(string)
	if !ok {
		return false, 0
	}

	depth = 0
	// Find all matches in the input string
	for {
		match := GetTemplatedString(strSource, dest)
		if !match {
			break
		}

		strSource = *dest
		depth++
	}
	if depth == 0 {
		return false, depth
	}

	*dest = strSource

	return true, depth
}

// recursively inspect src.
//   - e.g :we have src: {{ {{ {{ a }} }} }}.
//   - from : {a: b, b: c, c: d, d: { e; f; }}
//   - {{ a }} --> b ==> src: {{ {{ b }} }}.
//   - {{ b }} --> c ==> src: {{ c }}.
//   - {{ c }} --> d ==> will lookup to field d.
//   - it will get from key d then ' { e; f; }
func DeepTemplateEvaluateFromMap(mp map[string]interface{}, src interface{}, dest *interface{}) bool {
	keyMapField := ""

	// check if its using template
	if match, deep := DeepParseFromTemplate(src, &keyMapField); match {
		var (
			valueField interface{}
			lookUpMap  map[string]interface{} = mp
			ok         bool
		)

		for i := 0; i < deep; i++ {
			if valueField != nil {
				keyMapField, ok = valueField.(string)
				if !ok {
					break
				}
			}

			valueField, _ = LookUpMap(lookUpMap, keyMapField)
			if valueField == nil {

				return false
			}

		}

		*dest = valueField

		return true
	}

	return false
}
