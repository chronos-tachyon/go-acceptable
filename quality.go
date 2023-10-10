package acceptable

import (
	"encoding"
	"fmt"
	"math"
	"regexp"
	"strconv"
)

const kDigits = "0123456789"

var (
	reQuality  = regexp.MustCompile(`^(?:0|1|1\.0+|0?\.[0-9]+)$`)
	reQuality0 = regexp.MustCompile(`^(?:0|0\.0+|\.0+)$`)
	reQuality1 = regexp.MustCompile(`^1(?:\.0+)?$`)
)

type Quality uint

const (
	MinQuality = 0
	MaxQuality = 1000
)

func (q Quality) Append(out []byte) []byte {
	if q >= 1000 {
		return append(out, '1')
	}

	if q <= 0 {
		return append(out, '0')
	}

	a := kDigits[(q/100)%10]
	b := kDigits[(q/10)%10]
	c := kDigits[q%10]

	if c == '0' && b == '0' {
		return append(out, '0', '.', a)
	}

	if c == '0' {
		return append(out, '0', '.', a, b)
	}

	return append(out, '0', '.', a, b, c)
}

func (q Quality) GoString() string {
	return strconv.FormatUint(uint64(q), 10)
}

func (q Quality) String() string {
	var scratch [5]byte
	return string(q.Append(scratch[:0]))
}

func (q Quality) MarshalText() ([]byte, error) {
	return q.Append(make([]byte, 0, 5)), nil
}

func (q *Quality) Parse(input string) error {
	*q = 0

	if reQuality0.MatchString(input) {
		return nil
	}

	if reQuality1.MatchString(input) {
		*q = 1000
		return nil
	}

	if !reQuality.MatchString(input) {
		return fmt.Errorf("invalid quality %q", input)
	}

	f64, err := strconv.ParseFloat(input, 32)
	if err != nil {
		panic(err)
	}

	*q = Quality(math.RoundToEven(f64 * 1000))
	return nil
}

func (q *Quality) UnmarshalText(input []byte) error {
	return q.Parse(string(input))
}

func (q Quality) CompareTo(other Quality) int {
	return compareUints(uint(q), uint(other))
}

var (
	_ fmt.Stringer             = Quality(0)
	_ encoding.TextMarshaler   = Quality(0)
	_ encoding.TextUnmarshaler = (*Quality)(nil)
)
