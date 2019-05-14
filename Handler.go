package main

import (
	"models"
	"repositories"
)

// SalesOrder Structure for Sales Order
type SalesOrder struct {
	repo repositories.SalesOrder
}

// MKMReport Structure for MKM Report
type MKMReport struct {
	repo repositories.ReportMKM
}

// ReconReport Structure for MKM Report
type ReconReport struct {
	repo repositories.ReconData
}

// FeeConfig Structure for MKM Report
type FeeConfig struct {
	repo repositories.FeeConfig
}

// NewSOHandler Handler for SalesOrder
func NewSOHandler(db *models.DB) *SalesOrder {
	return &SalesOrder{
		repo: repositories.NewSalesOrder(db.SQL),
	}
}

// NewMKMHandler Handler for MKM report
func NewMKMHandler(db *models.DB) *MKMReport {
	return &MKMReport{
		repo: repositories.NewReportMKM(db.SQL),
	}
}

// NewReconHandler Handler for MKM report
func NewReconHandler(db *models.DB) *ReconReport {
	return &ReconReport{
		repo: repositories.NewRecon(db.SQL),
	}
}

// NewFeeConfigHandler Handler for MKM report
func NewFeeConfigHandler(db *models.DB) *FeeConfig {
	return &FeeConfig{
		repo: repositories.NewFeeConfig(db.SQL),
	}
}
