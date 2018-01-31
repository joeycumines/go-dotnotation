package dotnotation

import (
	"testing"
	"github.com/go-test/deep"
	"strings"
)

func TestGet(t *testing.T) {
	target := []interface{}{0, 1, 2}

	value, err := Get(target, "1")

	if nil != err || value != 1 {
		t.Fatalf("unexpected value %v / error %v", value, err)
	}
}

func TestSet(t *testing.T) {
	target := []interface{}{0, 1, 2}

	err := Set(target, "1", 3)

	if nil != err {
		t.Fatalf("unexpected error %v", err)
	}

	if diff := deep.Equal([]interface{}{0, 3, 2}, target); diff != nil {
		t.Fatalf("unexpected diff: %v", strings.Join(diff, ", "))
	}
}
