package acceptable

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestQuality_String(t *testing.T) {
	type testCase struct {
		Input  Quality
		Expect string
	}

	testData := [...]testCase{
		{0, "0"},
		{1, "0.001"},
		{9, "0.009"},
		{10, "0.01"},
		{90, "0.09"},
		{99, "0.099"},
		{100, "0.1"},
		{900, "0.9"},
		{990, "0.99"},
		{999, "0.999"},
		{1000, "1"},
	}

	for _, row := range testData {
		name := fmt.Sprintf("%d", uint(row.Input))
		t.Run(name, func(t *testing.T) {
			actual := row.Input.String()
			if actual != row.Expect {
				t.Errorf("wrong result:\n\texpect: %q\n\tactual: %q", row.Expect, actual)
			}
		})
	}
}

func TestQuality_Parse(t *testing.T) {
	type testCase struct {
		Input  string
		Expect Quality
		Err    error
	}

	testData := [...]testCase{
		{"0", 0, nil},
		{"1", 1000, nil},

		{"0.000", 0, nil},
		{"0.001", 1, nil},
		{"0.009", 9, nil},
		{"0.010", 10, nil},
		{"0.090", 90, nil},
		{"0.099", 99, nil},
		{"0.100", 100, nil},
		{"0.990", 990, nil},
		{"0.999", 999, nil},
		{"1.000", 1000, nil},

		{".000", 0, nil},
		{".001", 1, nil},
		{".009", 9, nil},
		{".010", 10, nil},
		{".090", 90, nil},
		{".099", 99, nil},
		{".100", 100, nil},
		{".990", 990, nil},
		{".999", 999, nil},

		{"", 0, errors.New("invalid quality \"\"")},
		{"2", 0, errors.New("invalid quality \"2\"")},
		{"0.", 0, errors.New("invalid quality \"0.\"")},
		{"1.", 0, errors.New("invalid quality \"1.\"")},
		{"1.1", 0, errors.New("invalid quality \"1.1\"")},
	}

	for _, row := range testData {
		t.Run(row.Input, func(t *testing.T) {
			var actual Quality
			err := actual.Parse(row.Input)
			if !reflect.DeepEqual(err, row.Err) {
				t.Errorf("wrong error:\n\texpect: %v\n\tactual: %v", row.Err, err)
			}
			if actual != row.Expect {
				t.Errorf("wrong result:\n\texpect: %v\n\tactual: %v", row.Expect, actual)
			}
		})
	}
}
