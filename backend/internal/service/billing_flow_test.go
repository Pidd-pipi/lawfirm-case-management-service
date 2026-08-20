package service

import (
	"testing"

	"cylawcase/internal/constants"
)

func TestBillingStatusCanFlowVoid(t *testing.T) {
	if billingStatusCanFlow(constants.BillingStatusVoid, constants.BillingStatusVoid) {
		t.Fatal("void -> void should be rejected")
	}
}

func TestBillingStatusCanFlowPaid(t *testing.T) {
	if billingStatusCanFlow(constants.BillingStatusInvoiced, constants.BillingStatusPaid) {
		t.Fatal("invoiced -> paid should be rejected")
	}
}

func TestBillingStatusCanFlowPending(t *testing.T) {
	if billingStatusCanFlow(constants.BillingStatusVoid, constants.BillingStatusPending) {
		t.Fatal("void -> pending should be rejected")
	}
}

func TestBillingStatusCanFlowInvoiced(t *testing.T) {
	if billingStatusCanFlow(constants.BillingStatusPending, constants.BillingStatusInvoiced) {
		t.Fatal("pending -> invoiced should be rejected")
	}
}
