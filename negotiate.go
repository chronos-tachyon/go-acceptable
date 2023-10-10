package acceptable

import (
	"bytes"
	"regexp"
	"sort"
	"strings"
	"sync"
)

var (
	gPatternMutex sync.Mutex
	gPatternCache map[string]*regexp.Regexp
)

const (
	kPatternCacheSize = 16
	kPatternCacheMax  = 32
)

func Negotiate(available, preferences List) (Acceptable, bool) {
	available = maybeSort(available)
	preferences = maybeSort(preferences)

	if available == nil {
		return Acceptable{}, false
	}

	if preferences == nil {
		var best Acceptable
		var hasBest bool
		for _, a := range available {
			if a.Quality > 0 && (!hasBest || a.Quality > best.Quality) {
				best = a
				hasBest = true
			}
		}
		return best, hasBest
	}

	list := make(candidateList, 0, len(available))
	for _, a := range available {
		var best Acceptable
		var hasBest bool

		for _, p := range preferences {
			if !isMatchingValue(a.Value, p.Value) {
				continue
			}

			if !isMatchingValue(a.SubValue, p.SubValue) {
				continue
			}

			if !isMatchingParams(a.Params, p.Params) {
				continue
			}

			if a.Quality <= 0 || p.Quality <= 0 {
				continue
			}

			best = p
			hasBest = true
			break
		}

		var q float64
		if hasBest {
			aq := float64(a.Quality) / 1000
			bq := float64(best.Quality) / 1000
			q = aq * bq
		}
		list = append(list, candidate{a, best, q})
	}

	list.Sort()
	return list[0].a, true
}

func maybeSort(list List) List {
	if len(list) <= 0 {
		return nil
	}
	if list.IsSorted() {
		return list
	}
	dupe := make(List, len(list))
	copy(dupe, list)
	dupe.Sort()
	return dupe
}

func isMatchingValue(actual, pattern string) bool {
	switch {
	case pattern == "":
		return actual == ""

	case pattern == "*":
		return actual != ""

	case strings.EqualFold(pattern, actual):
		return true

	default:
		return compilePattern(pattern).MatchString(actual)
	}
}

func isMatchingParams(actual, pattern map[string]string) bool {
	for key, pv := range pattern {
		av, found := actual[key]
		if !found || av != pv {
			return false
		}
	}
	return true
}

func compilePattern(pattern string) *regexp.Regexp {
	gPatternMutex.Lock()
	defer gPatternMutex.Unlock()

	if rx, found := gPatternCache[pattern]; found {
		return rx
	}

	rx := regexp.MustCompile(patternToRegexp(pattern))
	if gPatternCache == nil {
		gPatternCache = make(map[string]*regexp.Regexp, kPatternCacheSize)
	}
	for len(gPatternCache) >= kPatternCacheMax {
		var victim string
		for key := range gPatternCache {
			victim = key
			break
		}
		delete(gPatternCache, victim)
	}
	gPatternCache[pattern] = rx
	return rx
}

func patternToRegexp(pattern string) string {
	buf := gPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		gPool.Put(buf)
	}()

	buf.WriteString("(?i)^")
	for _, ch := range pattern {
		if ch == '*' {
			buf.WriteString(".+")
		} else {
			buf.WriteString(regexp.QuoteMeta(string(ch)))
		}
	}
	buf.WriteByte('$')
	return buf.String()
}

type candidate struct {
	a Acceptable
	p Acceptable
	q float64
}

type candidateList []candidate

func (list candidateList) Len() int {
	return len(list)
}

func (list candidateList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list candidateList) Less(i, j int) bool {
	a, b := list[i], list[j]
	cmp := compareFloats(a.q, b.q)
	if cmp != 0 {
		return cmp > 0
	}
	cmp = a.a.CompareTo(b.a)
	if cmp != 0 {
		return cmp < 0
	}
	return false
}

func (list candidateList) Sort() {
	sort.Sort(list)
}
