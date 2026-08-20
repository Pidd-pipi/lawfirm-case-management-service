package util

import (
	"reflect"
	"testing"
)

func TestGroupDigitsIsolated(t *testing.T) {
	a := groupDigits("1234567")
	_ = groupDigits("89")
	if !reflect.DeepEqual(a, []string{"1", "234", "567"}) {
		t.Fatalf("groupDigits 返回的切片被后续调用污染: %v", a)
	}
}

func TestGroupDigitsNegative(t *testing.T) {
	got := groupDigits("-1234")
	want := []string{"-", "1", "234"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("groupDigits(-1234)=%v, want %v", got, want)
	}
}
