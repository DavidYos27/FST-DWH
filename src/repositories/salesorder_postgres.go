package repositories

import (
	"database/sql"
	"fmt"
	"models"
	"strconv"
)

type postgresSalesOrder struct {
	Conn *sql.DB
}

// NewSalesOrder retunrs implement of post repository interface
func NewSalesOrder(Conn *sql.DB) SalesOrder {
	return &postgresSalesOrder{
		Conn: Conn,
	}
}

// GetSalesData get Data from Recharge Sales by ID
func (db *postgresSalesOrder) GetSalesData(startDate, endDate string) ([]*models.Order, error) {
	rows, err := db.Conn.Query("SELECT * FROM FACT.RECHARGE_SALES WHERE CREATED_DATE_ID BETWEEN $1 AND $2 ",
		startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		o := new(models.Order)
		err := rows.Scan(
			&o.RechargeID,
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
			&o.IsUsed,
		)
		if err != nil {
			return nil, err
		}
		tempToInt, _ := strconv.Atoi(o.SerialNumber)
		tempToString := fmt.Sprintf("%010v", tempToInt)
		o.SerialNumber = fmt.Sprintf(tempToString[len(tempToString)-10:])
		orders = append(orders, o)
	}
	return orders, err
}
