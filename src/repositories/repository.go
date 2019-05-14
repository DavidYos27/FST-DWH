package repositories

import "models"

// ReportMKM Interface for ReportMKM Model
type ReportMKM interface {
	GetMKMReport(startDate, endDate string) ([]*models.ReportMKM, error)
}

// SalesOrder Interface for Payfazz Model
type SalesOrder interface {
	GetSalesData(startDate, endDate string) ([]*models.Order, error)
}

// ReconData Interface for Recon Model
type ReconData interface {
	InsertToMatch(match []*models.Match) error
	InsertToNotMatch(notMatch []*models.NotMatch) error
}

// FeeConfig Interface for Fee
type FeeConfig interface {
	GetProviderConfig(providerID, productID string) (*models.Config, error)
}
