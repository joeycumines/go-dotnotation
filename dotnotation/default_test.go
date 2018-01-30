package dotnotation

import (
	"testing"
	"github.com/go-test/deep"
	"strings"
)

type dummyStruct struct {
	value string
}

type sliceGetterCase struct {
	name     string
	target   []interface{}
	property string
	success  bool
	result   interface{}
}

type sliceSetterCase struct {
	name     string
	target   []interface{}
	property string
	success  bool
	value    interface{}
	output   []interface{}
}

type mapGetterCase struct {
	name     string
	target   map[string]interface{}
	property string
	success  bool
	result   interface{}
}

type mapSetterCase struct {
	name     string
	target   map[string]interface{}
	property string
	success  bool
	value    interface{}
	output   map[string]interface{}
}

func TestDefaultGetter_invalidTypes(t *testing.T) {
	invalidTypes := []interface{}{
		1,
		"string",
		true,
		1.321321,
		nil,
		dummyStruct{"value"},
		[]dummyStruct{},
		map[string]dummyStruct{},
	}

	for _, value := range invalidTypes {
		v, err := DefaultGetter(value, "1")

		if err == nil {
			t.Errorf("expected error for %T %v", value, value)
		}

		if v != nil {
			t.Errorf("expected nil value for %T %v", value, value)
		}
	}
}

func TestDefaultGetter_slice(t *testing.T) {
	testCases := []sliceGetterCase{
		{
			name:     "non-integer property",
			target:   []interface{}{1, 2},
			property: "PROPERTY",
			success:  false,
			result:   nil,
		},
		{
			name:     "lower out of bounds",
			target:   []interface{}{1, 2},
			property: "-1",
			success:  false,
			result:   nil,
		},
		{
			name:     "upper out of bounds",
			target:   []interface{}{1, 2},
			property: "2",
			success:  false,
			result:   nil,
		},
		{
			name:     "upper",
			target:   []interface{}{1, 2},
			property: "1",
			success:  true,
			result:   2,
		},
		{
			name:     "lower",
			target:   []interface{}{1, 2},
			property: "0",
			success:  true,
			result:   1,
		},
	}

	for _, testCase := range testCases {
		v, err := DefaultGetter(testCase.target, testCase.property)

		if testCase.success {
			if err != nil || v != testCase.result {
				t.Errorf("%s failed: unexpected value %v / error %v for %v", testCase.name, v, err, testCase)
			}
		} else {
			if err == nil || v != nil {
				t.Errorf("%s failed: unexpected value %v / error %v for %v", testCase.name, v, err, testCase)
			}
		}
	}
}

func TestDefaultGetter_map(t *testing.T) {
	testCases := []mapGetterCase{
		{
			name: "missing property",
			target: map[string]interface{}{
				"one": 1,
				"two": 2,
			},
			property: "PROPERTY",
			success:  false,
			result:   nil,
		},
		{
			name: "match",
			target: map[string]interface{}{
				"one": 1,
				"two": 2,
			},
			property: "two",
			success:  true,
			result:   2,
		},
	}

	for _, testCase := range testCases {
		v, err := DefaultGetter(testCase.target, testCase.property)

		if testCase.success {
			if err != nil || v != testCase.result {
				t.Errorf("%s failed: unexpected value %v / error %v for %v", testCase.name, v, err, testCase)
			}
		} else {
			if err == nil || v != nil {
				t.Errorf("%s failed: unexpected value %v / error %v for %v", testCase.name, v, err, testCase)
			}
		}
	}
}

func TestDefaultSetter_invalidTypes(t *testing.T) {
	invalidTypes := []interface{}{
		1,
		"string",
		true,
		1.321321,
		nil,
		dummyStruct{"value"},
		[]dummyStruct{},
		map[string]dummyStruct{},
	}

	for _, value := range invalidTypes {
		copied := value

		err := DefaultSetter(&value, "1", 1)

		if err == nil {
			t.Errorf("expected error for %T %v", value, value)
		}

		// this doesn't actually do anything for the struct etc lol
		if diff := deep.Equal(copied, value); diff != nil {
			t.Errorf("expected value to be unchanged, but was %v", diff)
		}
	}
}

func TestDefaultSetter_slice(t *testing.T) {
	testCases := []sliceSetterCase{
		{
			name:     "non-integer property",
			target:   []interface{}{1, 2},
			property: "PROPERTY",
			success:  false,
			value:    3,
			output:   []interface{}{1, 2},
		},
		{
			name:     "lower out of bounds",
			target:   []interface{}{1, 2},
			property: "-1",
			success:  false,
			value:    3,
			output:   []interface{}{1, 2},
		},
		{
			name:     "upper out of bounds",
			target:   []interface{}{1, 2},
			property: "3",
			success:  false,
			value:    3,
			output:   []interface{}{1, 2},
		},
		{
			name:     "upper",
			target:   []interface{}{1, 2},
			property: "1",
			success:  true,
			value:    3,
			output:   []interface{}{1, 3},
		},
		{
			name:     "lower",
			target:   []interface{}{1, 2},
			property: "0",
			success:  true,
			value:    3,
			output:   []interface{}{3, 2},
		},
		{
			name:     "append",
			target:   []interface{}{1, 2},
			property: "2",
			success:  true,
			value:    3,
			output:   []interface{}{1, 2, 3},
		},
	}

	for _, testCase := range testCases {
		err := DefaultSetter(&(testCase.target), testCase.property, testCase.value)

		if diff := deep.Equal(testCase.output, testCase.target); diff != nil {
			t.Errorf("%s failed: unexpected diff (%v) for %v", testCase.name, strings.Join(diff, ", "), testCase)
		}

		if testCase.success {
			if err != nil {
				t.Errorf("%s failed: unexpected error %v for %v", testCase.name, err, testCase)
			}
		} else {
			if err == nil {
				t.Errorf("%s failed: unexpected error %v for %v", testCase.name, err, testCase)
			}
		}
	}
}

func TestDefaultSetter_map(t *testing.T) {
	testCases := []mapSetterCase{
		{
			name: "miss",
			target: map[string]interface{}{
				"one": 1,
				"two": 2,
			},
			property: "PROPERTY",
			success:  true,
			value:    3,
			output: map[string]interface{}{
				"one":      1,
				"two":      2,
				"PROPERTY": 3,
			},
		},
		{
			name: "collision",
			target: map[string]interface{}{
				"one": 1,
				"two": 2,
			},
			property: "two",
			success:  true,
			value:    3,
			output: map[string]interface{}{
				"one": 1,
				"two": 3,
			},
		},
	}

	for _, testCase := range testCases {
		err := DefaultSetter(&(testCase.target), testCase.property, testCase.value)

		if diff := deep.Equal(testCase.output, testCase.target); diff != nil {
			t.Errorf("%s failed: unexpected diff (%v) for %v", testCase.name, strings.Join(diff, ", "), testCase)
		}

		if testCase.success {
			if err != nil {
				t.Errorf("%s failed: unexpected error %v for %v", testCase.name, err, testCase)
			}
		} else {
			if err == nil {
				t.Errorf("%s failed: unexpected error %v for %v", testCase.name, err, testCase)
			}
		}
	}
}
