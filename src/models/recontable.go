package models

// Match Structure for Recon Data
type Match struct {
	MatchID    int
	ProviderID string
	ProductID  string
	Source     string
	PzID       string
	BsID       string
	Amount     int
	Status     int
	DateID     int
	FeeID      int
	CreatedAt  string
	CreatedBy  string
	UpdatedAt  string
	UpdatedBy  string
}

// NotMatch Structure for Recon Data
type NotMatch struct {
	ProviderID string
	ProductID  string
	Source     string
	KeyID      string
	Amount     int
	Status     int
	DateID     int
	IsUsed     int
	CreatedAt  string
	CreatedBy  string
	UpdatedAt  string
	UpdatedBy  string
}
