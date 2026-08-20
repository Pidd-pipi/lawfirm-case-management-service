package dto

import "testing"

func TestPageQueryNormalizeNilReceiver(t *testing.T) {
	var p *PageQuery
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Normalize nil receiver panicked: %v", r)
		}
	}()
	p.Normalize()
}
