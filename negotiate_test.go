package acceptable

import (
	"reflect"
	"testing"
)

func TestNegotiate(t *testing.T) {
	type testCase struct {
		Name        string
		Available   List
		Preferences List
		Expect      Acceptable
		ExpectOK    bool
	}

	testData := [...]testCase{
		{
			Name: "Empty",
		},
		{
			Name: "NoneAvailable",
			Preferences: List{
				{"text", "html", nil, 1000},
				{"text", "*", nil, 900},
				{"*", "*", nil, 100},
			},
		},
		{
			Name: "NonePreferred",
			Available: List{
				{"text", "html", nil, 1000},
				{"application", "json", nil, 999},
			},
			Expect:   Acceptable{"text", "html", nil, 1000},
			ExpectOK: true,
		},
		{
			Name: "Specificity",
			Available: List{
				{"text", "html", nil, 1000},
				{"text", "plain", paramsCharset, 1000},
			},
			Preferences: List{
				{"text", "*", paramsCharset, 1000},
				{"text", "html", nil, 999},
			},
			Expect:   Acceptable{"text", "plain", paramsCharset, 1000},
			ExpectOK: true,
		},
	}

	for _, row := range testData {
		t.Run(row.Name, func(t *testing.T) {
			actual, ok := Negotiate(row.Available, row.Preferences)
			if ok != row.ExpectOK || !reflect.DeepEqual(actual, row.Expect) {
				t.Errorf("wrong result:\n\texpect: %#v, %t\n\tactual: %#v, %t", row.Expect, row.ExpectOK, actual, ok)
			}
		})
	}
}
