package acceptable

import (
	"math"
	"sort"
)

func compareUints(a, b uint) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func compareFloats(a, b float64) int {
	if math.IsNaN(a) || math.IsNaN(b) {
		panic("number is NaN")
	}
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func compareStrings(a, b string) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func compareValues(a, b string) int {
	switch {
	case a == b:
		return 0
	case a == "*":
		return 1
	case b == "*":
		return -1
	case a < b:
		return -1
	default:
		return 1
	}
}

func compareParams(a, b map[string]string) int {
	if cmp := compareUints(uint(len(a)), uint(len(b))); cmp != 0 {
		return -cmp
	}

	aKeys, bKeys := paramKeys(a), paramKeys(b)
	for i := range aKeys {
		aKey, bKey := aKeys[i], bKeys[i]
		if cmp := compareStrings(aKey, bKey); cmp != 0 {
			return cmp
		}

		aValue, bValue := a[aKey], b[bKey]
		if cmp := compareStrings(aValue, bValue); cmp != 0 {
			return cmp
		}
	}

	return 0
}

func paramKeys(params map[string]string) []string {
	if len(params) <= 0 {
		return nil
	}

	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
