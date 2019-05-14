package repositories

import (
	"database/sql"
	"models"
)

type postgresConfig struct {
	Conn *sql.DB
}

// NewFeeConfig return implement of post repository interface
func NewFeeConfig(Conn *sql.DB) FeeConfig {
	return &postgresConfig{
		Conn: Conn,
	}
}

// GetProviderConfig get Config for specific provider
func (db *postgresConfig) GetProviderConfig(providerID, productID string) (*models.Config, error) {
	rows, err := db.Conn.Query("SELECT * FROM DIMENSION.PROVIDER_CONFIG WHERE PROVIDERID = $1 AND PRODUCTID = $2",
		providerID, productID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	config := new(models.Config)
	for rows.Next() {

		err := rows.Scan(&config.ProviderID,
			&config.ProductID,
			&config.Discount,
			&config.Fee,
			&config.Coupon,
		)
		if err != nil {
			return nil, err
		}
	}
	return config, err
}
