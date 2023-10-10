package acceptable

import (
	"reflect"
	"testing"
)

func TestList_String(t *testing.T) {
	type testCase struct {
		Name   string
		Input  List
		Expect string
	}

	testData := [...]testCase{
		{
			Name:   "Empty",
			Input:  nil,
			Expect: "",
		},
		{
			Name: "One",
			Input: List{
				{"text", "html", nil, 1000},
			},
			Expect: "text/html",
		},
		{
			Name: "Three",
			Input: List{
				{"text", "html", nil, 1000},
				{"text", "*", nil, 900},
				{"*", "*", nil, 100},
			},
			Expect: "text/html, text/*;q=0.9, */*;q=0.1",
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

func TestList_Parse(t *testing.T) {
	type testCase struct {
		Name   string
		Input  string
		Mode   SubValueMode
		Expect List
		Err    error
	}

	testData := [...]testCase{
		{
			Name:   "Empty",
			Input:  "",
			Expect: nil,
		},
		{
			Name:  "One",
			Input: "text/html",
			Expect: List{
				{"text", "html", nil, 1000},
			},
		},
		{
			Name:  "Three",
			Input: "text/html, text/*;q=0.9, */*;q=0.1",
			Expect: List{
				{"text", "html", nil, 1000},
				{"text", "*", nil, 900},
				{"*", "*", nil, 100},
			},
		},
	}

	for _, row := range testData {
		t.Run(row.Name, func(t *testing.T) {
			var actual List
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

func TestList_Sort(t *testing.T) {
	type testCase struct {
		Name   string
		Input  List
		Expect List
	}

	testData := [...]testCase{
		{
			Name:   "Empty",
			Input:  nil,
			Expect: nil,
		},
		{
			Name: "One",
			Input: List{
				{"text", "html", nil, 1000},
			},
			Expect: List{
				{"text", "html", nil, 1000},
			},
		},
		{
			Name: "Three",
			Input: List{
				{"*", "*", nil, 100},
				{"text", "*", nil, 900},
				{"text", "html", nil, 1000},
			},
			Expect: List{
				{"text", "html", nil, 1000},
				{"text", "*", nil, 900},
				{"*", "*", nil, 100},
			},
		},
		{
			Name: "ThreeSameQ",
			Input: List{
				{"*", "*", nil, 1000},
				{"text", "*", nil, 1000},
				{"text", "*", paramsCharset, 1000},
				{"text", "html", nil, 1000},
				{"text", "html", paramsCharset, 1000},
			},
			Expect: List{
				{"text", "html", paramsCharset, 1000},
				{"text", "*", paramsCharset, 1000},
				{"text", "html", nil, 1000},
				{"text", "*", nil, 1000},
				{"*", "*", nil, 1000},
			},
		},
		{
			Name: "ThreeSameValue",
			Input: List{
				{"text", "html", nil, 1000},
				{"text", "html", nil, 900},
				{"text", "html", nil, 100},
			},
			Expect: List{
				{"text", "html", nil, 100},
				{"text", "html", nil, 900},
				{"text", "html", nil, 1000},
			},
		},
		{
			Name: "ComplexOne",
			Input: List{
				{"*", "*", nil, 1000},
				{"text", "html", nil, 1000},
				{"text", "*", nil, 1000},
				{"image", "*", nil, 1000},
				{"image", "webp", nil, 0},
				{"application", "xhtml+xml", nil, 1000},
			},
			Expect: List{
				{"application", "xhtml+xml", nil, 1000},
				{"image", "webp", nil, 0},
				{"image", "*", nil, 1000},
				{"text", "html", nil, 1000},
				{"text", "*", nil, 1000},
				{"*", "*", nil, 1000},
			},
		},
		{
			Name: "ComplexTwo",
			Input: List{
				{"*", "*", nil, 1000},
				{"text", "html", nil, 1000},
				{"text", "*", paramsCharset, 1000},
				{"image", "*", nil, 1000},
				{"image", "webp", nil, 0},
				{"application", "xhtml+xml", nil, 1000},
			},
			Expect: List{
				{"text", "*", paramsCharset, 1000},
				{"application", "xhtml+xml", nil, 1000},
				{"image", "webp", nil, 0},
				{"image", "*", nil, 1000},
				{"text", "html", nil, 1000},
				{"*", "*", nil, 1000},
			},
		},
	}
	for _, row := range testData {
		t.Run(row.Name, func(t *testing.T) {
			actual := row.Input
			actual.Sort()
			if !reflect.DeepEqual(actual, row.Expect) {
				t.Errorf("wrong result:\n\texpect: %v\n\tactual: %v", row.Expect, actual)
			}
		})
	}
}
