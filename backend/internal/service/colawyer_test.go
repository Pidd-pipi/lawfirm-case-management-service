package service

import (
	"encoding/json"
	"testing"
)

func TestJSONCoLawyersDedupes(t *testing.T) {
	raw := jsonCoLawyers([]uint64{2, 2, 3})
	var got []uint64
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != 2 || got[1] != 3 {
		t.Fatalf("jsonCoLawyers = %v, want [2 3]", got)
	}
}

func TestJSONCoLawyersDropsZero(t *testing.T) {
	raw := jsonCoLawyers([]uint64{0, 2})
	var got []uint64
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != 2 {
		t.Fatalf("jsonCoLawyers = %v, want [2]", got)
	}
}
