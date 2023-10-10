package acceptable

import (
	"encoding"
	"fmt"
	"strings"
)

type SubValueMode uint

const (
	OptionalSubValue SubValueMode = iota
	RequiredSubValue
	AbsentSubValue
)

type Acceptable struct {
	Value    string
	SubValue string
	Params   map[string]string
	Quality  Quality
}

func (a Acceptable) Append(out []byte) []byte {
	if a.Value == "" {
		return out
	}

	out = appendToken(out, a.Value)

	if a.SubValue != "" {
		out = append(out, '/')
		out = appendToken(out, a.SubValue)
	}

	if n := len(a.Params); n > 0 {
		keys := paramKeys(a.Params)
		for _, key := range keys {
			out = append(out, ';')
			out = appendToken(out, key)
			out = append(out, '=')
			out = appendToken(out, a.Params[key])
		}
	}

	if a.Quality < 1000 {
		out = append(out, ";q="...)
		out = a.Quality.Append(out)
	}

	return out
}

func (a Acceptable) String() string {
	return string(a.Append(nil))
}

func (a Acceptable) MarshalText() ([]byte, error) {
	return a.Append(nil), nil
}

func (a *Acceptable) Parse(input string, mode SubValueMode) error {
	*a = Acceptable{}

	input = consumeSpace(input)

	value, rest, ok := consumeToken(input)
	if !ok {
		return fmt.Errorf("expect token, got %q", input)
	}
	input = rest

	input = consumeSpace(input)

	var subValue string
	var hasSubValue bool

	switch mode {
	case OptionalSubValue:
		hasSubValue = strings.HasPrefix(input, "/")

	case RequiredSubValue:
		hasSubValue = true
		if !strings.HasPrefix(input, "/") {
			return fmt.Errorf("expect '/', got %q", input)
		}

	case AbsentSubValue:
		hasSubValue = false

	default:
		return fmt.Errorf("unknown mode %v", mode)
	}

	if hasSubValue {
		input = input[1:]
		input = consumeSpace(input)

		subValue, rest, ok = consumeToken(input)
		if !ok {
			return fmt.Errorf("expect token, got %q", input)
		}
		input = rest

		input = consumeSpace(input)
	}

	var p map[string]string
	var q Quality = 1000
	for strings.HasPrefix(input, ";") {
		input = input[1:]
		input = consumeSpace(input)

		var paramName string
		paramName, rest, ok = consumeToken(input)
		if !ok {
			return fmt.Errorf("expect token, got %q", input)
		}
		input = rest

		input = consumeSpace(input)
		if !strings.HasPrefix(input, "=") {
			return fmt.Errorf("expect '=', got %q", input)
		}
		input = input[1:]
		input = consumeSpace(input)

		var paramValue string
		paramValue, rest, ok = consumeQuoted(input)
		if !ok {
			return fmt.Errorf("expect token or quoted string, got %q", input)
		}
		input = rest

		input = consumeSpace(input)

		paramName = strings.ToLower(paramName)
		if paramName == "q" {
			if err := q.Parse(paramValue); err != nil {
				return fmt.Errorf("%w", err)
			}
		} else {
			if p == nil {
				p = make(map[string]string)
			}
			p[paramName] = paramValue
		}
	}

	if input != "" {
		return fmt.Errorf("expect ';', got %q", input)
	}

	*a = Acceptable{value, subValue, p, q}
	return nil
}

func (a *Acceptable) UnmarshalText(input []byte) error {
	return a.Parse(string(input), OptionalSubValue)
}

func (a Acceptable) EqualTo(b Acceptable) bool {
	return a.CompareTo(b) == 0
}

func (a Acceptable) CompareTo(b Acceptable) int {
	if cmp := compareParams(a.Params, b.Params); cmp != 0 {
		return cmp
	}
	if cmp := compareValues(a.Value, b.Value); cmp != 0 {
		return cmp
	}
	if cmp := compareValues(a.SubValue, b.SubValue); cmp != 0 {
		return cmp
	}
	if cmp := a.Quality.CompareTo(b.Quality); cmp != 0 {
		return cmp
	}
	return 0
}

var (
	_ fmt.Stringer             = Acceptable{}
	_ encoding.TextMarshaler   = Acceptable{}
	_ encoding.TextUnmarshaler = (*Acceptable)(nil)
)
