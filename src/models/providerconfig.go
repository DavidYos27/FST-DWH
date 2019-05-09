package models

type Config struct{
	ProviderID 	string
	ProductID 	string
	Discount	string
	Fee			int
	Coupon		string
}

// GetProviderConfig get Config for specific provider
func (db *DB) GetProviderConfig(providerID,productID string)(*Config,error){
	rows, err := db.Query("SELECT * FROM DIMENSION.PROVIDER_CONFIG WHERE PROVIDERID = $1 AND PRODUCTID = $2",
	providerID, productID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	config := new(Config)
	for rows.Next() {
		
		err := rows.Scan( 	&config.ProviderID,
							&config.ProductID,
							&config.Discount,
							&config.Fee,
							&config.Coupon,
		)
		if err != nil {
			return nil, err
		}
	}
	return config,err
}
