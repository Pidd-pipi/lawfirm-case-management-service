package constants

import "testing"

func TestIsValidBillingStatusInvoiced(t *testing.T) {
	if !IsValidBillingStatus(BillingStatusInvoiced) {
		t.Fatal("IsValidBillingStatus(invoiced) should be true")
	}
}
