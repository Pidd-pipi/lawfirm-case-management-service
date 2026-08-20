package constants

import "testing"

func TestBillingStatusTransitions(t *testing.T) {
	m := BillingStatusTransitions()
	if m == nil {
		t.Fatal("BillingStatusTransitions should return non-nil map")
	}
	if !m[BillingStatusPaid][BillingStatusInvoiced] {
		t.Fatal("paid -> invoiced should be allowed")
	}
}
