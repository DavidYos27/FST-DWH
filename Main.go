package main

import (
	"fmt"
	"log"
	"models"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const dwhFaConn = "host=10.0.42.77 port=5433 user=davin.timothy password=jklgfsfa64 dbname=dwh_fa_operation sslmode=disable"
const dbPaymentConn = "host=10.0.89.85 port=6432 user=dwh-app password=68610b0ab7d8484ad7f126227b834633 dbname=payment sslmode=disable"

// Env Structure for Provider Model
type Env struct {
	db models.ProviderModel
}

var notRecon = make([]*models.NotMatch, 0)
var recon = make([]*models.Match, 0)

const (
	providerID     = "1"
	productID      = "1"
	providerSource = "1"
	pzSource       = "2"
	notMatchStatus = -1
	matchStatus    = 1
)

func main() {
	// startSalesReport()

	db, err := models.NewDB(dwhFaConn)

	if err != nil {
		log.Panic(err)
	}

	env := &Env{db}
	start := "20190101"
	end := "20190114"
	totalNotRecon := 0

	config, err := env.db.GetProviderConfig(providerID, productID)
	orders, err := env.db.GetSalesData(start, end)

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

	reports, err := env.db.GetMKMReport(start, end)
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
							notMatchStatus,
							array.CreatedDateID,
							1)
						totalNotRecon += array.SellPrice - array.Fee
						log.Println(array)
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
							notMatchStatus,
							reportsMap[key][i].CreatedDateID,
							2)
						totalNotRecon += reportsMap[key][i].SellPrice - reportsMap[key][i].Fee
						log.Println(reportsMap[key][i])
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
						notMatchStatus,
						array.CreatedDateID, 3)
					totalNotRecon += array.SellPrice - array.Fee
					log.Println(array)
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
					notMatchStatus,
					array.CreatedDateID,
					4)
				totalNotRecon += array.SellPrice - array.Fee
				log.Println(array)
			}
		}
	}

	starts := time.Now()

	errors := env.db.InsertToMatch(recon)
	elapsed := time.Since(starts)
	log.Printf("Inserting took %s second", elapsed)
}

// InsertRecon Insert data into corresponding map
func InsertRecon(status bool, providerID, productID, source, pzID, bsID string, match, dateID, from int) {
	if status == false {
		data := new(models.NotMatch)
		data.ProviderID = providerID
		data.ProductID = productID
		data.Source = source
		data.KeyID = pzID
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
