package dotnotation

import (
	"testing"
	"github.com/go-test/deep"
	"strings"
)

type accessorTestCase struct {
	name    string
	key     string
	target  interface{}
	before  interface{}
	getErr  bool
	after   interface{}
	setErr  bool
	outcome interface{}
}

func TestAccessor(t *testing.T) {
	accessors := []Accessor{
		// default accessor
		{},
		{
			Getter: DefaultGetter,
			Setter: DefaultSetter,
			Parser: DefaultParser,
		},
	}

	for _, accessor := range accessors {
		testCases := []accessorTestCase{
			{
				name:   "empty map",
				key:    "one",
				target: map[string]interface{}{},
				before: nil,
				getErr: true,
				after:  1,
				setErr: false,
				outcome: map[string]interface{}{
					"one": 1,
				},
			},
			{
				name:    "empty map nested set failure",
				key:     "one.two",
				target:  map[string]interface{}{},
				before:  nil,
				getErr:  true,
				after:   1,
				setErr:  true,
				outcome: map[string]interface{}{},
			},
			{
				name: "existing map set collision",
				key:  "one",
				target: map[string]interface{}{
					"one": 1,
					"two": 2,
				},
				before: 1,
				getErr: false,
				after:  3,
				setErr: false,
				outcome: map[string]interface{}{
					"one": 3,
					"two": 2,
				},
			},
			{
				name: "existing map set miss",
				key:  "three",
				target: map[string]interface{}{
					"one": 1,
					"two": 2,
				},
				before: nil,
				getErr: true,
				after:  3,
				setErr: false,
				outcome: map[string]interface{}{
					"one":   1,
					"two":   2,
					"three": 3,
				},
			},
			{
				name: "existing map set collision on nil value",
				key:  "three",
				target: map[string]interface{}{
					"one":   1,
					"two":   2,
					"three": nil,
				},
				before: nil,
				getErr: false,
				after:  3,
				setErr: false,
				outcome: map[string]interface{}{
					"one":   1,
					"two":   2,
					"three": 3,
				},
			},
			{
				name: "existing map set collision nested",
				key:  "lvl1.lvl2.one",
				target: map[string]interface{}{
					"lvl1": map[string]interface{}{
						"lvl2": map[string]interface{}{
							"one": 1,
							"two": 2,
						},
					},
				},
				before: 1,
				getErr: false,
				after:  3,
				setErr: false,
				outcome: map[string]interface{}{
					"lvl1": map[string]interface{}{
						"lvl2": map[string]interface{}{
							"one": 3,
							"two": 2,
						},
					},
				},
			},
		}

		for _, testCase := range testCases {
			// check before
			actual, err := accessor.Get(testCase.target, testCase.key)

			if err != nil && !testCase.getErr { // allow on get error
				t.Errorf("%s failed: %v", testCase.name, err)
				continue
			}

			if diff := deep.Equal(testCase.before, actual); diff != nil {
				t.Errorf("%s failed: unexpected diff %v", testCase.name, strings.Join(diff, ", "))
				continue
			}

			// set
			err = accessor.Set(testCase.target, testCase.key, testCase.after)

			if err != nil && !testCase.setErr { // allow on set error
				t.Errorf("%s failed: %v", testCase.name, err)
				continue
			}

			// check after
			actual, err = accessor.Get(testCase.target, testCase.key)

			if err != nil && !(testCase.getErr && testCase.setErr) { // only allow error on both get and set error
				t.Errorf("%s failed: %v", testCase.name, err)
				continue
			}

			// the actual comparison for the after diff will be the before if we expected a set error
			if !testCase.setErr {
				if diff := deep.Equal(testCase.after, actual); diff != nil {
					t.Errorf("%s failed: unexpected diff %v", testCase.name, strings.Join(diff, ", "))
					continue
				}
			} else {
				if diff := deep.Equal(testCase.before, actual); diff != nil {
					t.Errorf("%s failed: unexpected diff %v", testCase.name, strings.Join(diff, ", "))
					continue
				}
			}

			// check outcome
			if diff := deep.Equal(testCase.outcome, testCase.target); diff != nil {
				t.Errorf("%s failed: unexpected diff %v", testCase.name, strings.Join(diff, ", "))
				continue
			}
		}
	}
}

func TestAccessor_Get_parserError(t *testing.T) {
	accessor := Accessor{
		Parser: func(key string) []string {
			return nil
		},
	}

	_, err := accessor.Get(nil, "key")

	if "no properties parsed from key: key" != err.Error() {
		t.Fatalf("unexpected error %v", err)
	}
}
func TestAccessor_Set_parserError(t *testing.T) {
	accessor := Accessor{
		Parser: func(key string) []string {
			return nil
		},
	}

	err := accessor.Set(nil, "key", nil)

	if "no properties parsed from key: key" != err.Error() {
		t.Fatalf("unexpected error %v", err)
	}
}
