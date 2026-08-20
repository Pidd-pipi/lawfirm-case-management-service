package util

import (
	"reflect"
	"testing"
)

func TestFormatMoneyListIsolated(t *testing.T) {
	a := FormatMoneyList([]float64{1.5, 2.25})
	_ = FormatMoneyList([]float64{9.99})
	if !reflect.DeepEqual(a, []string{"¥1.50", "¥2.25"}) {
		t.Fatalf("FormatMoneyList 返回的切片被后续调用污染: %v", a)
	}
}

func TestFormatMoneyListContent(t *testing.T) {
	got := FormatMoneyList([]float64{0, 1234.5})
	want := []string{"¥0.00", "¥1234.50"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FormatMoneyList = %v, want %v", got, want)
	}
}
