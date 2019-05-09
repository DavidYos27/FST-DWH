package models
import (
	// "fmt"
	// "log"
)
// ReportMKM Structure for Model Provider Report MKM
type ReportMKM struct{
	ID				string
	ProviderID 		int
	SerialNumber 	string
	SellPrice		int
	Fee				int
	UploadTime		string
	CreatedDateID	int
}

// GetMKMReport get Data from Provider Report MKM
func (db *DB) GetMKMReport(startDate, endDate string) ([]*ReportMKM, error) {
	rows, err := db.Query("SELECT * FROM FACT.PROVIDER_REPORT_MKM WHERE CREATED_DATE_ID BETWEEN $1 AND $2 ",
				startDate, endDate)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []*ReportMKM

	for rows.Next() {
		r := new(ReportMKM)
		err := rows.Scan( 	&r.ID,
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
		_,r.SellPrice = Abs(r.SellPrice)
		// r.SerialNumber = fmt.Sprintf("%v|%v",r.SerialNumber,r.SellPrice)
		reports = append(reports, r)
	}
	return reports,err
}

// Abs Return check value is input minus or not and return Absolute value
func Abs( integer int) (bool,int){
	isMinus := true
	if integer < 0{
		integer = integer * -1
		isMinus = false
	}
	return isMinus, integer
}