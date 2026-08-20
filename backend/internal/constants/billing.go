package constants

// BillingType 费用类型枚举。
const (
	BillingTypeAttorneyFee = "attorney_fee"
	BillingTypeCourtFee    = "court_fee"
	BillingTypeTravelFee   = "travel_fee"
	BillingTypeOther       = "other"
)

// BillingTypeValues 全部费用类型值。
var BillingTypeValues = []string{BillingTypeAttorneyFee, BillingTypeCourtFee, BillingTypeTravelFee, BillingTypeOther}

// BillingStatus 账单状态枚举。
const (
	BillingStatusPending  = "pending"
	BillingStatusPaid     = "paid"
	BillingStatusInvoiced = "invoiced"
	BillingStatusVoid     = "void"
)

// BillingStatusValues 全部账单状态值。
var BillingStatusValues = []string{BillingStatusPending, BillingStatusPaid, BillingStatusInvoiced, BillingStatusVoid}

// IsValidBillingType 校验费用类型。
func IsValidBillingType(s string) bool {
	for _, v := range BillingTypeValues {
		if v == s {
			return true
		}
	}
	return false
}

// IsValidBillingStatus 校验账单状态。
func IsValidBillingStatus(s string) bool {
	for _, v := range BillingStatusValues {
		if v == s {
			return true
		}
	}
	return false
}

// BillingStatusTransitions 返回账单状态机迁移表。
// 状态流转：pending -> paid -> invoiced；任意非 void 状态均可 -> void；void 为终态，不可再流转。
// 注意：未支付（pending）的账单不能直接开票，必须先支付（paid）。
func BillingStatusTransitions() map[string]map[string]bool {
	return map[string]map[string]bool{
		BillingStatusPending: {
			BillingStatusPaid: true,
			BillingStatusVoid: true,
		},
		BillingStatusPaid: {
			BillingStatusInvoiced: true,
			BillingStatusVoid:    true,
		},
		BillingStatusInvoiced: {
			BillingStatusVoid: true,
		},
		BillingStatusVoid: {},
	}
}

// DocumentFileType 文档类型枚举。
const (
	DocTypeComplaint = "complaint"
	DocTypeDefense   = "defense"
	DocTypeEvidence  = "evidence"
	DocTypeJudgment  = "judgment"
	DocTypeContract  = "contract"
	DocTypeOther     = "other"
)

// DocumentFileTypeValues 全部文档类型值。
var DocumentFileTypeValues = []string{DocTypeComplaint, DocTypeDefense, DocTypeEvidence, DocTypeJudgment, DocTypeContract, DocTypeOther}
