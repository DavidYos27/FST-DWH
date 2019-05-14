package repositories

import (
	"database/sql"
	"models"
)

type postgresReportMKM struct {
	Conn *sql.DB
}

// NewReportMKM retunrs implement of post repository interface
func NewReportMKM(Conn *sql.DB) ReportMKM {
	return &postgresReportMKM{
		Conn: Conn,
	}
}

// GetMKMReport get Data from Provider Report MKM
func (db *postgresReportMKM) GetMKMReport(startDate, endDate string) ([]*models.ReportMKM, error) {
	rows, err := db.Conn.Query("SELECT * FROM FACT.PROVIDER_REPORT_MKM WHERE CREATED_DATE_ID BETWEEN $1 AND $2",
		startDate, endDate)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*models.ReportMKM

	for rows.Next() {
		r := new(models.ReportMKM)
		err := rows.Scan(&r.ID,
			&r.ProviderID,
			&r.SerialNumber,
			&r.SellPrice,
			&r.Fee,
			&r.UploadTime,
			&r.CreatedDateID,
		)
		if err != nil {
			return nil, err
		}
		_, r.SellPrice = Abs(r.SellPrice)
		reports = append(reports, r)
	}
	return reports, err
}

// Abs Return check value is input minus or not and return Absolute value
func Abs(integer int) (bool, int) {
	isMinus := true
	if integer < 0 {
		integer = integer * -1
		isMinus = false
	}
	return isMinus, integer
}
