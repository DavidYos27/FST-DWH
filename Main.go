package main

import (
	"log"
	"models"

	_ "github.com/lib/pq"
)

const dwhFaConn = "host=10.0.42.77 port=5433 user=davin.timothy password=jklgfsfa64 dbname=dwh_fa_operation sslmode=disable"
const dbPaymentConn = "host=10.0.89.85 port=6432 user=dwh-app password=68610b0ab7d8484ad7f126227b834633 dbname=payment sslmode=disable"

func main() {
	// startSalesReport()
	// startMkmReport()
	db, err := models.NewDB(dwhFaConn)
	if err != nil {
		log.Panic(err)
	}

	start := "20190401"
	end := "20190416"
	totalNotRecon := 0

	StartReconMKM(start, end, totalNotRecon, db)

}
