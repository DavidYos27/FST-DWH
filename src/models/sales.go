package models

import (
	"strconv"
	"fmt"
	// "log"
)

// Order Structure for Model Sales
type Order struct{
	RechargeID 		string
	OrderID 		int
	ProductID 		int
	ProviderID 		int
	OperatorID 		int
	PlanID 			int
	SellPrice 		int
	InvoiceNetPrice 		int
	Amount			int
	NetPrice		int
	Fee 			int
	Netfee 			int
	Discount 		int
	CouponID 		int
	CouponAmount 	int
	RefNo 			string
	ItemID 			string
	SerialNumber 	string
	Refunded	 	int
	RefundedAmount	int
	Status 			string
	Cashback 		int
	CashbackMembershipID int
	UserID 			int
	CreatedAt 		string 
	UpdatedAt 		string
	PostedAt 		string
	CreatedDateID 	int
	UpdatedDateID 	int
	PostedDateID 	int
	isUsed			int
}

// GetSalesData get Data from Recharge Sales by ID
func (db *DB) GetSalesData(startDate, endDate string) ([]*Order, error) {
	rows, err := db.Query("SELECT * FROM FACT.RECHARGE_SALES WHERE CREATED_DATE_ID BETWEEN $1 AND $2 ",
				startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*Order
	
	for rows.Next() {
		o := new(Order)
		err := rows.Scan(	&o.RechargeID,
							&o.OrderID,
							&o.ProductID,
							&o.ProviderID,
							&o.OperatorID,
							&o.PlanID,
							&o.SellPrice,
							&o.InvoiceNetPrice,
							&o.Amount,
							&o.NetPrice,
							&o.Fee,
							&o.Netfee,
							&o.Discount,
							&o.CouponID,
							&o.CouponAmount,
							&o.RefNo,
							&o.ItemID,
							&o.SerialNumber,
							&o.Refunded,
							&o.RefundedAmount,
							&o.Status,
							&o.Cashback,
							&o.CashbackMembershipID,
							&o.UserID,
							&o.CreatedAt,
							&o.UpdatedAt,
							&o.PostedAt,
							&o.CreatedDateID,
							&o.UpdatedDateID,
							&o.PostedDateID,
							&o.isUsed,		
			)
		if err != nil {			
			return nil, err
		}

		tempToInt,_ 	:= strconv.Atoi(o.SerialNumber)
		tempToString   	:= fmt.Sprintf("%010v", tempToInt)
		o.SerialNumber 	 = fmt.Sprintf(tempToString[len(tempToString)-10:])
		orders = append(orders, o)
	}
	return orders,err
}