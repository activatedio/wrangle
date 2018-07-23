package template

import (
	"reflect"
	"testing"
)

func Test_join(t *testing.T) {

	cases := map[string]struct {
		in              []interface{}
		sepprefixsuffix []string
		expected        string
	}{
		"empty": {
			sepprefixsuffix: []string{","},
		},
		"simple": {
			in:              []interface{}{"1", "2", "3"},
			sepprefixsuffix: []string{","},
			expected:        "1,2,3",
		},
		"simple-prefix-suffixed": {
			in:              []interface{}{"1", "2", "3"},
			sepprefixsuffix: []string{",", "\"", "\""},
			expected:        `"1","2","3"`,
		},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {

			got := join(v.in, v.sepprefixsuffix...)

			if got != v.expected {

				t.Fatalf("Wanted [%s] got [%s]", v.expected, got)
			}
		})
	}

}

type simpleStruct struct {
	A string
}

type simple2Struct struct {
	A string
	B string
}

func Test_project(t *testing.T) {

	cases := map[string]struct {
		input    []interface{}
		property string
		expected []interface{}
	}{
		"single-map": {
			[]interface{}{
				map[string]interface{}{
					"a": "a1",
				},
			},
			"a",
			[]interface{}{"a1"},
		},
		"single-struct": {
			[]interface{}{
				simpleStruct{
					A: "a1",
				},
			},
			"A",
			[]interface{}{"a1"},
		},
		"simple-struct-ptr": {
			[]interface{}{
				&simpleStruct{
					A: "a1",
				},
			},
			"A",
			[]interface{}{"a1"},
		},
		"multi-map": {
			[]interface{}{
				map[string]interface{}{
					"a": "a1",
					"b": "b1",
				},
				map[string]interface{}{
					"a": "a2",
					"b": "b2",
				},
			},
			"b",
			[]interface{}{"b1", "b2"},
		},
		"multi-struct": {
			[]interface{}{
				simple2Struct{
					A: "a1",
					B: "b1",
				},
				simple2Struct{
					A: "a2",
					B: "b2",
				},
			},
			"B",
			[]interface{}{"b1", "b2"},
		},
		"multi-struct-ptr": {
			[]interface{}{
				&simple2Struct{
					A: "a1",
					B: "b1",
				},
				&simple2Struct{
					A: "a1",
					B: "b2",
				},
			},
			"B",
			[]interface{}{"b1", "b2"},
		},
	}

	for k, v := range cases {

		t.Run(k, func(t *testing.T) {
			got := project(v.input, v.property)
			if !reflect.DeepEqual(got, v.expected) {
				t.Fatalf("Expected %+v, got %+v", v.expected, got)
			}
		})
	}
}
