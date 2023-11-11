package common

func MergeMap(src, dest map[string]interface{}) map[string]interface{} {
	for k, v := range dest {
		(src)[k] = v
	}
	return src
}
