package util

import (
	"reflect"
	"testing"
)

func TestUniquePreservesInput(t *testing.T) {
	input := []uint64{1, 2, 1, 3}
	orig := append([]uint64(nil), input...)
	got := Unique(input)
	if !reflect.DeepEqual(got, []uint64{1, 2, 3}) {
		t.Fatalf("Unique = %v, want [1 2 3]", got)
	}
	if !reflect.DeepEqual(input, orig) {
		t.Fatalf("Unique 修改了入参: %v -> %v", orig, input)
	}
}

func TestFilterPreservesInput(t *testing.T) {
	input := []uint64{1, 0, 2, 0, 3}
	orig := append([]uint64(nil), input...)
	got := Filter(input, func(v uint64) bool { return v > 0 })
	if !reflect.DeepEqual(got, []uint64{1, 2, 3}) {
		t.Fatalf("Filter = %v, want [1 2 3]", got)
	}
	if !reflect.DeepEqual(input, orig) {
		t.Fatalf("Filter 修改了入参: %v -> %v", orig, input)
	}
}
