package models

// Order Structure for Model Sales
type Order struct {
	RechargeID           string
	OrderID              int
	ProductID            int
	ProviderID           int
	OperatorID           int
	PlanID               int
	SellPrice            int
	InvoiceNetPrice      int
	Amount               int
	NetPrice             int
	Fee                  int
	Netfee               int
	Discount             int
	CouponID             int
	CouponAmount         int
	RefNo                string
	ItemID               string
	SerialNumber         string
	Refunded             int
	RefundedAmount       int
	Status               string
	Cashback             int
	CashbackMembershipID int
	UserID               int
	CreatedAt            string
	UpdatedAt            string
	PostedAt             string
	CreatedDateID        int
	UpdatedDateID        int
	PostedDateID         int
	IsUsed               int
}
