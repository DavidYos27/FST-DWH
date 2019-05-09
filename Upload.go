package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"models"
	"os"
	"strconv"
	"time"
)

const (
	gopher = 10
)

func startMkmReport() {
	for i := 1; i <= 31; i++ {
		month := "01"
		year := "2019"
		date, _ := strconv.Atoi(fmt.Sprintf("%v%v%02d", year, month, i))
		day := fmt.Sprintf("Upload/MKM/%02d-%v-%v.csv", i, month, year)
		println(day, date)
		UploadDataMKM(day, date)
	}
}

func startSalesReport() {
	UploadDataSales("Upload/Sales/Recharge MKM 2019-01.csv")
}

// UploadDataSales for Uploading data sales to DB
func UploadDataSales(path string) {
	var entries int
	var sales []models.Order
	counter := 0
	csvFile, _ := os.Open(path)
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil &&
			counter > 0 {
		}

		if counter > 0 {
			sales = append(sales, models.Order{
				RechargeID:           line[0],
				OrderID:              StringtoInt(line[1]),
				ProductID:            StringtoInt(line[2]),
				ProviderID:           StringtoInt(line[3]),
				OperatorID:           StringtoInt(line[4]),
				PlanID:               StringtoInt(line[5]),
				SellPrice:            StringtoInt(line[6]),
				InvoiceNetPrice:      StringtoInt(line[7]),
				Amount:               StringtoInt(line[8]),
				NetPrice:             StringtoInt(line[9]),
				Fee:                  StringtoInt(line[10]),
				Netfee:               StringtoInt(line[11]),
				Discount:             StringtoInt(line[12]),
				CouponID:             StringtoInt(line[13]),
				CouponAmount:         StringtoInt(line[14]),
				RefNo:                line[15],
				ItemID:               line[16],
				SerialNumber:         line[17],
				Refunded:             StringtoInt(line[18]),
				RefundedAmount:       StringtoInt(line[19]),
				Status:               line[20],
				Cashback:             StringtoInt(line[21]),
				CashbackMembershipID: StringtoInt(line[22]),
				UserID:               StringtoInt(line[23]),
				CreatedAt:            line[24],
				UpdatedAt:            line[25],
				PostedAt:             line[26],
				CreatedDateID:        StringtoInt(line[27]),
				UpdatedDateID:        StringtoInt(line[28]),
				PostedDateID:         StringtoInt(line[29]),
			})
		}
		counter++
	}

	finishChan := make(chan int)
	entries = (len(sales) / gopher) + 1
	sStmt := "INSERT INTO FACT.RECHARGE_SALES(RECHARGE_ID,ORDER_ID,PRODUCT_ID,PROVIDER_ID,OPERATOR_ID,PLAN_ID,INVOICE_AMOUNT,INVOICE_NETPRICE,AMOUNT,NETPRICE,FEE,NETFEE,DISCOUNT,COUPON_ID,COUPON_AMOUNT,REF_NO,ITEM_ID,SERIAL_NUMBER,REFUNDED,REFUNDED_AMOUNT,STATUS,CASHBACK,CASHBACK_MEMBERSHIP_ID,USER_ID,CREATED_AT,UPDATED_AT,POSTED_AT,CREATED_DATE_ID,UPDATED_DATE_ID,POSTED_DATE_ID) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30)"

	for i := 0; i < gopher; i++ {
		// log.Println("insert:", i, entries)
		go insertDataSales(i, entries, sStmt, finishChan, sales)
	}

	finishedGophers := 0
	finishLoop := false

	for {
		if finishLoop {
			break
		}
		select {
		case n := <-finishChan:
			finishedGophers += n
			if finishedGophers == gopher {
				finishLoop = true
				println("Finished")
			}
		}
	}
}

func StringtoInt(strings string) int {
	result, _ := strconv.Atoi(strings)
	return result
}

func insertDataSales(gopher, entries int, sStmt string, c chan int, reports []models.Order) {
	db, err := sql.Open("postgres", dwhFaConn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	stmt, err := db.Prepare(sStmt)
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()
	for i := (gopher * entries); i < (gopher*entries)+entries; i++ {
		if i < len(reports) {
			res, err := stmt.Exec(
				reports[i].RechargeID,
				reports[i].OrderID,
				reports[i].ProductID,
				reports[i].ProviderID,
				reports[i].OperatorID,
				reports[i].PlanID,
				reports[i].SellPrice,
				reports[i].InvoiceNetPrice,
				reports[i].Amount,
				reports[i].NetPrice,
				reports[i].Fee,
				reports[i].Netfee,
				reports[i].Discount,
				reports[i].CouponID,
				reports[i].CouponAmount,
				reports[i].RefNo,
				reports[i].ItemID,
				reports[i].SerialNumber,
				reports[i].Refunded,
				reports[i].RefundedAmount,
				reports[i].Status,
				reports[i].Cashback,
				reports[i].CashbackMembershipID,
				reports[i].UserID,
				reports[i].CreatedAt,
				reports[i].UpdatedAt,
				reports[i].PostedAt,
				reports[i].CreatedDateID,
				reports[i].UpdatedDateID,
				reports[i].PostedDateID)
			if err != nil || res == nil {
				log.Fatal(err)
			}
		} else {
			break
		}
	}
	c <- 1
}

// UploadDataMKM function to UploadReportMKM to DB
func UploadDataMKM(path string, date int) {
	var serialNumber string
	var entries int
	var reports []models.ReportMKM
	csvFile, _ := os.Open(path)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	counter := 0
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil &&
			counter > 0 {
			// log.Println(error)
			// log.Fatal(error)
		}

		if counter > 0 {
			sellPrice, _ := strconv.Atoi(line[2])
			fee, _ := strconv.Atoi(line[3])
			times := time.Now().String()

			if len(line[5]) >= 11 {
				serialNumber = line[5][len(line[5])-10 : len(line[5])]
			} else {
				serialNumber = ""
			}
			reports = append(reports, models.ReportMKM{
				ProviderID:    1,
				SerialNumber:  serialNumber,
				SellPrice:     sellPrice,
				Fee:           fee,
				UploadTime:    times,
				CreatedDateID: date,
			})
		}
		counter++
	}

	finishChan := make(chan int)
	entries = (len(reports) / gopher) + 1
	sStmt := "INSERT INTO FACT.PROVIDER_REPORT_MKM(PROVIDER_ID,SERIAL_NUMBER,SELL_PRICE,FEE,UPLOAD_TIME,CREATED_DATE_ID) VALUES ($1, $2, $3, $4, $5, $6)"

	for i := 0; i < gopher; i++ {
		// log.Println("insert:", i, entries)
		go insertData(i, entries, sStmt, finishChan, reports)
	}

	finishedGophers := 0
	finishLoop := false

	for {
		if finishLoop {
			break
		}
		select {
		case n := <-finishChan:
			finishedGophers += n
			if finishedGophers == gopher {
				finishLoop = true
				println("Finished")
			}
		}
	}
}

func insertData(gopher, entries int, sStmt string, c chan int, reports []models.ReportMKM) {
	db, err := sql.Open("postgres", dwhFaConn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	stmt, err := db.Prepare(sStmt)
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()
	for i := (gopher * entries); i < (gopher*entries)+entries; i++ {
		if i < len(reports) {
			res, err := stmt.Exec(
				reports[i].ProviderID,
				reports[i].SerialNumber,
				reports[i].SellPrice,
				reports[i].Fee,
				time.Now().String(),
				reports[i].CreatedDateID)
			if err != nil || res == nil {
				log.Fatal(err)
			}
		} else {
			break
		}
	}
	c <- 1
}
