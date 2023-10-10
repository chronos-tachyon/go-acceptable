package acceptable

import (
	"encoding"
	"fmt"
	"sort"
)

type List []Acceptable

func (list List) Append(out []byte) []byte {
	for i, a := range list {
		if i > 0 {
			out = append(out, ", "...)
		}
		out = a.Append(out)
	}
	return out
}

func (list List) String() string {
	return string(list.Append(nil))
}

func (list List) MarshalText() ([]byte, error) {
	return list.Append(nil), nil
}

func (list *List) Parse(input string, mode SubValueMode) error {
	*list = nil

	var result List

	type pstate uint
	const (
		rootState pstate = iota
		quoteState
		escapeState
	)

	state := rootState
	limit := uint(len(input))
	start := uint(0)
	for i := uint(0); i < limit; i++ {
		ch := input[i]
		switch {
		case state == escapeState:
			state = quoteState
		case state == quoteState && ch == '\\':
			state = escapeState
		case state == quoteState && ch == '"':
			state = rootState
		case state == rootState && ch == '"':
			state = quoteState
		case state == rootState && ch == ',':
			var a Acceptable
			if str := consumeSpace(input[start:i]); str != "" {
				if err := a.Parse(str, mode); err != nil {
					return err
				}
				result = append(result, a)
			}
			start = i + 1
		}
	}

	if str := consumeSpace(input[start:]); str != "" {
		var a Acceptable
		if err := a.Parse(str, mode); err != nil {
			return err
		}
		result = append(result, a)
	}

	*list = result
	return nil
}

func (list *List) UnmarshalText(input []byte) error {
	return list.Parse(string(input), OptionalSubValue)
}

func (list List) Len() int {
	return len(list)
}

func (list List) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list List) Less(i, j int) bool {
	a, b := list[i], list[j]
	cmp := a.CompareTo(b)
	return cmp < 0
}

func (list List) Sort() {
	sort.Sort(list)
}

func (list List) IsSorted() bool {
	return sort.IsSorted(list)
}

var (
	_ fmt.Stringer             = List(nil)
	_ encoding.TextMarshaler   = List(nil)
	_ encoding.TextUnmarshaler = (*List)(nil)
)
