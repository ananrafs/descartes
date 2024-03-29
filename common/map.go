package common

import (
	"strings"
)

type mapManipulator func(map[string]interface{})

func ManipulateMap(srcs ...map[string]interface{}) mapManipulator {
	return func(m map[string]interface{}) {
		for _, src := range srcs {
			for k, v := range src {
				m[k] = v
			}
		}
	}
}

func (mm mapManipulator) ChainCopy(src map[string]interface{}) mapManipulator {
	return func(m map[string]interface{}) {
		mm(m)
		mm.Copy(src)
	}
}

func (mm mapManipulator) Copy(src map[string]interface{}) (dest map[string]interface{}) {
	dest = make(map[string]interface{})
	mm(dest)
	for k, v := range src {
		dest[k] = v
	}
	return
}

func (mm mapManipulator) DeepCopy(src map[string]interface{}) (dest map[string]interface{}) {
	d := make(map[string]interface{})
	mm(d)
	return deepCopyMap(d, src)
}

func (mm mapManipulator) Merge(src map[string]interface{}) (dest map[string]interface{}) {
	dest = make(map[string]interface{})
	mm(dest)
	for k, v := range src {
		(dest)[k] = v
	}
	return dest
}

func deepCopyMap(init, src map[string]interface{}) (dest map[string]interface{}) {
	dest = init
	for key, val := range src {
		var copiedValue = val
		childMap, ok := val.(map[string]interface{})
		if ok {
			_nested := deepCopyMap(make(map[string]interface{}), childMap)
			copiedValue = _nested
		}

		dest[key] = copiedValue
	}

	return
}

func LookUpMap(mp map[string]interface{}, src string) (interface{}, error) {
	return recursivelyLookupMap(mp, 0, strings.Split(src, "."))
}

func recursivelyLookupMap(mp map[string]interface{}, index int, source []string) (res interface{}, err error) {
	if index >= len(source) {
		return nil, nil
	}
	if val, ok := mp[source[index]]; ok {
		childMap, ok := val.(map[string]interface{})
		if !ok {
			return val, nil
		}
		index++
		if index >= len(source) {
			return val, nil
		}
		lookup, err := recursivelyLookupMap(childMap, index, source)
		if err != nil {
			return nil, err
		}
		return lookup, nil
	}
	return nil, ErrorNotFoundOnMap(source[index])
}

// ExtractMap is method to evaluate `source` interface{} and fill it to dest
//
//	if source is map, then it will recursively lookup to key-pairs.
//	modifiers is used to alter or modify value of map.
//	e.g: lookup a.b.c --> wil search `c` in : a { b: { c: {d: 12, e: 13} }}
//	and will set `*dest` to {d: 12, e: 13}
func ExtractMap(source interface{}, dest *map[string]interface{}, modifiers ...func(*string, *interface{})) bool {
	sMap, ok := source.(map[string]interface{})
	if !ok {
		return false
	}

	copiedMap := ManipulateMap(map[string]interface{}{}).Copy(sMap)
	for key, val := range copiedMap {
		var _dest map[string]interface{}
		isObj := ExtractMap(val, &_dest, modifiers...)
		if isObj {
			copiedMap[key] = _dest
			continue
		}
		for _, mod := range modifiers {
			mod(&key, &val)
		}
		copiedMap[key] = val
	}
	*dest = copiedMap

	return true
}

// SetMap is func to set key from given map a value
func SetMap(mp map[string]interface{}, key string, value interface{}) {
	objArr := strings.Split(key, ".")
	recursivelySetMap(mp, value, 0, objArr...)

}

// Recursive Set Map from given Keys.
//
//	e.g: keys : a.b.c.d; value :12.
//	will set a { b: { c: {d: 12} } }.
//	will modify struct if a/b/c is not struct/map.
//	will create nested struct if key not found
func recursivelySetMap(mp map[string]interface{}, value interface{}, index int, keys ...string) {
	if index >= len(keys) {
		return
	}
	var (
		childMap   map[string]interface{}
		currentKey = keys[index]
	)

	currValue, ok := mp[currentKey]
	if ok {
		ok := false
		childMap, ok = currValue.(map[string]interface{})
		if !ok {
			childMap = make(map[string]interface{})
		}

	} else {
		childMap = make(map[string]interface{})
	}

	index++
	if index >= len(keys) {
		mp[currentKey] = value
		return
	}

	mp[currentKey] = childMap
	recursivelySetMap(childMap, value, index, keys...)
}
