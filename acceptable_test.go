package acceptable

import (
	"fmt"
	"reflect"
	"testing"
)

var (
	paramsCharset = map[string]string{"charset": "utf-8"}
	paramsWeird   = map[string]string{"weird": "some text with \" and \\"}
)

func TestAcceptable_String(t *testing.T) {
	type testCase struct {
		Name   string
		Input  Acceptable
		Expect string
	}

	testData := [...]testCase{
		{
			Name:   "Empty",
			Input:  Acceptable{"", "", nil, 0},
			Expect: "",
		},
		{
			Name:   "ValueQ0",
			Input:  Acceptable{"foo", "", nil, 0},
			Expect: "foo;q=0",
		},
		{
			Name:   "ValueQ1",
			Input:  Acceptable{"foo", "", nil, 1000},
			Expect: "foo",
		},
		{
			Name:   "StarQ0",
			Input:  Acceptable{"*", "", nil, 0},
			Expect: "*;q=0",
		},
		{
			Name:   "StarQ1",
			Input:  Acceptable{"*", "", nil, 1000},
			Expect: "*",
		},
		{
			Name:   "ValueSubQ0",
			Input:  Acceptable{"foo", "bar", nil, 0},
			Expect: "foo/bar;q=0",
		},
		{
			Name:   "ValueSubQ0",
			Input:  Acceptable{"foo", "bar", nil, 1000},
			Expect: "foo/bar",
		},
		{
			Name:   "ValueStarQ0",
			Input:  Acceptable{"foo", "*", nil, 0},
			Expect: "foo/*;q=0",
		},
		{
			Name:   "ValueStarQ0",
			Input:  Acceptable{"foo", "*", nil, 1000},
			Expect: "foo/*",
		},
		{
			Name:   "StarStarQ0",
			Input:  Acceptable{"*", "*", nil, 0},
			Expect: "*/*;q=0",
		},
		{
			Name:   "StarStarQ0",
			Input:  Acceptable{"*", "*", nil, 1000},
			Expect: "*/*",
		},
		{
			Name:   "Text-HTML-UTF8",
			Input:  Acceptable{"text", "html", paramsCharset, 1000},
			Expect: "text/html;charset=utf-8",
		},
		{
			Name:   "Weird",
			Input:  Acceptable{"foo", "", paramsWeird, 1000},
			Expect: `foo;weird="some text with \" and \\"`,
		},
	}

	for _, row := range testData {
		t.Run(row.Name, func(t *testing.T) {
			actual := row.Input.String()
			if actual != row.Expect {
				t.Errorf("wrong result:\b\texpect: %q\n\tactual: %q", row.Expect, actual)
			}
		})
	}
}

func TestAcceptable_Parse(t *testing.T) {
	type testCase struct {
		Name   string
		Input  string
		Mode   SubValueMode
		Expect Acceptable
		Err    error
	}

	testData := [...]testCase{
		{
			Name:   "GZip",
			Input:  "gzip",
			Expect: Acceptable{"gzip", "", nil, 1000},
		},
		{
			Name:   "GZipQ1",
			Input:  "gzip;q=1",
			Expect: Acceptable{"gzip", "", nil, 1000},
		},
		{
			Name:   "GZipQ0",
			Input:  "gzip;q=0",
			Expect: Acceptable{"gzip", "", nil, 0},
		},
		{
			Name:   "Text-HTML",
			Input:  "text/html",
			Expect: Acceptable{"text", "html", nil, 1000},
		},
		{
			Name:   "Text-HTML-UTF8",
			Input:  "text/html;charset=utf-8",
			Expect: Acceptable{"text", "html", paramsCharset, 1000},
		},
		{
			Name:   "Text-HTML-UTF8-q0.1",
			Input:  "text/html;charset=utf-8;q=0.1",
			Expect: Acceptable{"text", "html", paramsCharset, 100},
		},
		{
			Name:   "Spaces",
			Input:  " text / html ; charset = utf-8 ; q = 0.1 ",
			Expect: Acceptable{"text", "html", paramsCharset, 100},
		},
		{
			Name:   "Weird",
			Input:  `foo;weird="some text with \" and \\"`,
			Expect: Acceptable{"foo", "", paramsWeird, 1000},
		},

		{
			Name:  "FailEmpty",
			Input: "",
			Err:   fmt.Errorf("expect token, got \"\""),
		},
		{
			Name:  "FailComma",
			Input: ",",
			Err:   fmt.Errorf("expect token, got \",\""),
		},
		{
			Name:  "FailSlash",
			Input: "/",
			Err:   fmt.Errorf("expect token, got \"/\""),
		},
		{
			Name:  "FailMissingSubValue",
			Input: "gzip;q=1",
			Mode:  RequiredSubValue,
			Err:   fmt.Errorf("expect '/', got \";q=1\""),
		},
		{
			Name:  "FailExtraSubValue",
			Input: "text/html;q=1",
			Mode:  AbsentSubValue,
			Err:   fmt.Errorf("expect ';', got \"/html;q=1\""),
		},
	}

	for _, row := range testData {
		t.Run(row.Name, func(t *testing.T) {
			var actual Acceptable
			err := actual.Parse(row.Input, row.Mode)
			if !reflect.DeepEqual(err, row.Err) {
				t.Errorf("wrong error:\n\texpect: %v\n\tactual: %v", row.Err, err)
			}
			if !reflect.DeepEqual(actual, row.Expect) {
				t.Errorf("wrong result:\n\texpect: %v\n\tactual: %v", row.Expect, actual)
			}
		})
	}
}
