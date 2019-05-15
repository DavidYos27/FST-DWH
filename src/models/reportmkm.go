package models

// ReportMKM Structure for Model Provider Report MKM
type ReportMKM struct {
	ID            string
	ProviderID    int
	SerialNumber  string
	SellPrice     int
	Fee           int
	UploadTime    string
	CreatedDateID int
}
