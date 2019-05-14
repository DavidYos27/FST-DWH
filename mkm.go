package main

import (
	"fmt"
	"log"
	"models"
	"strconv"
	"strings"
	"time"
)

const (
	providerID     = "1"
	productID      = "1"
	providerSource = "1"
	pzSource       = "5"
	notMatchStatus = -1
	matchStatus    = 1
)

var notRecon = make([]*models.NotMatch, 0)
var recon = make([]*models.Match, 0)

// StartReconMKM Start to Recon Data
func StartReconMKM(start, end string, totalNotRecon int, db *models.DB) {
	soHandler := NewSOHandler(db)
	mkmHandler := NewMKMHandler(db)
	reconHandler := NewReconHandler(db)
	configHandler := NewFeeConfigHandler(db)

	config, err := configHandler.repo.GetProviderConfig(providerID, productID)
	orders, err := soHandler.repo.GetSalesData(start, end)

	salesMap := make(map[string][]*models.Order)
	if err != nil {
		fmt.Printf("Error: %v", err)
	} else {
		for _, value := range orders {
			if value.Refunded != 1 {
				value.SellPrice = value.Amount + (value.Fee * config.Fee)
				key := fmt.Sprintf("%v|%v", value.SerialNumber, value.SellPrice)
				salesMap[key] = append(salesMap[key], value)
			}
		}
	}

	reports, err := mkmHandler.repo.GetMKMReport(start, end)
	reportsMap := make(map[string][]*models.ReportMKM)

	if err != nil {
		fmt.Printf("Error: %v", err)
	} else {
		for _, value := range reports {
			key := fmt.Sprintf("%v|%v", value.SerialNumber, value.SellPrice)
			reportsMap[key] = append(reportsMap[key], value)
		}
	}

	for key, value := range salesMap {
		if _, ok := reportsMap[key]; ok {
			for i, array := range value {
				if array.Refunded != 1 &&
					array.RechargeID != "" {
					if len(reportsMap[key]) >= i {
						if CheckStatus(array.RechargeID,
							reportsMap[key][i].ID,
							array.CreatedDateID,
							reportsMap[key][i].CreatedDateID); true {
							InsertRecon(true,
								providerID,
								productID,
								pzSource,
								array.RechargeID,
								reportsMap[key][i].ID,
								"",
								array.SellPrice,
								matchStatus,
								array.CreatedDateID,
								0)
						} else {
							InsertRecon(
								false,
								providerID,
								productID,
								pzSource,
								array.RechargeID,
								"",
								array.SerialNumber,
								array.SellPrice,
								notMatchStatus,
								array.CreatedDateID,
								1)
							totalNotRecon += array.SellPrice - array.Fee
						}
						array.RechargeID = ""
					} else {
						InsertRecon(
							false,
							providerID,
							productID,
							providerSource,
							array.RechargeID,
							"",
							array.SerialNumber,
							array.SellPrice,
							notMatchStatus,
							array.CreatedDateID,
							1)
						totalNotRecon += array.SellPrice - array.Fee
					}
				}
			}
			_, variant := Abs(len(reportsMap[key]) - len(value))
			if len(reportsMap[key]) > len(value) {
				for i := len(reportsMap[key]) - variant; i < len(reportsMap[key]); i++ {
					if isNumeric(reportsMap[key][i].SerialNumber) {
						InsertRecon(false,
							providerID,
							productID,
							providerSource,
							reportsMap[key][i].ID,
							"",
							reportsMap[key][i].SerialNumber,
							reportsMap[key][i].SellPrice,
							notMatchStatus,
							reportsMap[key][i].CreatedDateID,
							2)
						totalNotRecon += reportsMap[key][i].SellPrice - reportsMap[key][i].Fee
					}
				}
			}
			delete(reportsMap, key)
		} else {
			for _, array := range value {
				if array.Refunded != 1 {
					InsertRecon(false,
						providerID,
						productID,
						pzSource,
						array.RechargeID,
						"",
						array.SerialNumber,
						array.SellPrice,
						notMatchStatus,
						array.CreatedDateID, 3)
					totalNotRecon += array.SellPrice - array.Fee
				}
			}
		}
	}

	// Loop all Provider report without Payfazz data
	for _, value := range reportsMap {
		for _, array := range value {
			if isNumeric(array.SerialNumber) {
				InsertRecon(false,
					providerID,
					productID,
					providerSource,
					array.ID,
					"",
					array.SerialNumber,
					array.SellPrice,
					notMatchStatus,
					array.CreatedDateID,
					4)
				totalNotRecon += array.SellPrice - array.Fee
			}
		}
	}

	starts := time.Now()
	errors := reconHandler.repo.InsertToMatch(recon)
	if errors == nil {
		log.Println("insert Not match")
		errors = reconHandler.repo.InsertToNotMatch(notRecon)
	}
	if errors != nil {
		log.Println(errors)
	}
	elapsed := time.Since(starts)
	log.Printf("Inserting took %s second", elapsed)
	log.Println("length:", len(notRecon), totalNotRecon)

}

// InsertRecon Insert data into corresponding map
func InsertRecon(status bool, providerID, productID, source, pzID, bsID, serialNumber string, amount int, match, dateID, from int) {
	if status == false {
		data := new(models.NotMatch)
		data.ProviderID = providerID
		data.ProductID = productID
		data.Source = source
		data.KeyID = pzID
		data.Amount = amount
		data.Status = match
		data.DateID = dateID
		notRecon = append(notRecon, data)
	} else if status == true {
		data := new(models.Match)
		data.ProviderID = providerID
		data.ProductID = productID
		data.Source = source
		data.PzID = pzID
		data.BsID = bsID
		data.Amount = amount
		data.Status = match
		data.DateID = dateID
		recon = append(recon, data)
	}
}

// CheckStatus Check that reconciliation status
func CheckStatus(pzValue, proValue string, pzDate, proDate int) bool {
	tempDate := strconv.Itoa(pzDate)
	layout := "2006-01-02"
	pzConvDate, _ := time.Parse(layout, fmt.Sprintf("%v-%v-%v", tempDate[0:4], tempDate[4:6], tempDate[6:8]))
	pzConvDate = pzConvDate.AddDate(0, 0, -1)
	if pzDate == proDate ||
		proDate == TtoInt(pzConvDate, "-") {
		return true
	}
	return false
}

// isNumeric check whether the input is numeric or not
func isNumeric(text string) bool {
	if _, err := strconv.Atoi(text); err != nil {
		return false
	}
	return true
}

// TtoInt format time format to int
func TtoInt(time time.Time, delimiter string) int {
	result, _ := strconv.Atoi(strings.Replace(time.String(), delimiter, "", 2)[0:8])
	return result
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
