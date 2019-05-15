package models

import (
	"database/sql"
)

// Operator Structure for Table
type Operator struct {
	OperatorID   int
	OperatorName string
	CreatedAt    string
	UpdatedAt    sql.NullString
}

// GetOperatorID Get Provider ID from Provider Table
// func (db *DB) GetOperatorID() ([]*Operator, error) {
// 	rows, err := db.Query("SELECT * FROM DIMENSION.OPERATOR")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	operators := make([]*Operator, 0)

// 	for rows.Next() {
// 		o := new(Operator)
// 		err := rows.Scan(&o.OperatorID, &o.OperatorName, &o.CreatedAt, &o.UpdatedAt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		operators = append(operators, o)
// 	}
// 	return operators, nil
// }
