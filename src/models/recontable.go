package models

import (
	"log"
	"strconv"
)

// Match Structure for Recon Data
type Match struct {
	MatchID    int
	ProviderID string
	ProductID  string
	Source     string
	PzID       string
	BsID       string
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
	Status     int
	DateID     int
	IsUsed     int
	CreatedAt  string
	CreatedBy  string
	UpdatedAt  string
	UpdatedBy  string
}

const (
	entries = 500
)

// InsertToMatch insert to Match table
func (db *DB) InsertToMatch(tags []*Match) error {

	numFields := 6
	stmt := "INSERT INTO fact.report_match (provider_id, product_id, pz_id, bs_id, status,date_id) values"
	for i := 0; i <= len(tags)/entries; i++ {
		query := stmt
		values := []interface{}{}
		counter := 0
		if i == len(tags)/entries {
			for j := i * entries; j < len(tags); j++ {
				values = append(values, AToInt(tags[j].ProviderID), AToInt(tags[j].ProductID), tags[j].PzID, tags[j].BsID, tags[j].Status, tags[j].DateID)
				n := counter * numFields

				query += "("
				for k := 0; k < numFields; k++ {
					query += `$` + strconv.Itoa(n+k+1) + ","
				}
				query = query[:len(query)-1] + "),"
				counter++
			}
		} else {
			for j := i * entries; j < (i*entries)+entries; j++ {
				values = append(values, AToInt(tags[j].ProviderID), AToInt(tags[j].ProductID), tags[j].PzID, tags[j].BsID, tags[j].Status, tags[j].DateID)
				n := counter * numFields

				query += "("
				for k := 0; k < numFields; k++ {
					query += `$` + strconv.Itoa(n+k+1) + ","
				}
				query = query[:len(query)-1] + "),"
				counter++
			}
		}
		query = query[:len(query)-1]
		_, err := db.Exec(query, values...)
		if err != nil {
			log.Println(query)
			return err
		}
	}
	return nil
}

// InsertToNotMatch Insert data to Not Match Table
func (db *DB) InsertToNotMatch(tags []*NotMatch) error {
	numFields := 6
	stmt := "INSERT INTO fact.report_match (provider_id, product_id, pz_id, bs_id, status,date_id) values"
	for i := 0; i <= len(tags)/entries; i++ {
		query := stmt
		values := []interface{}{}
		counter := 0
		if i == len(tags)/entries {
			for j := i * entries; j < len(tags); j++ {
				values = append(values, AToInt(tags[j].ProviderID), AToInt(tags[j].ProductID), tags[j].KeyID, tags[j].Status, tags[j].DateID)
				n := counter * numFields

				query += "("
				for k := 0; k < numFields; k++ {
					query += `$` + strconv.Itoa(n+k+1) + ","
				}
				query = query[:len(query)-1] + "),"
				counter++
			}
		} else {
			for j := i * entries; j < (i*entries)+entries; j++ {
				values = append(values, AToInt(tags[j].ProviderID), AToInt(tags[j].ProductID), tags[j].PzID, tags[j].BsID, tags[j].Status, tags[j].DateID)
				n := counter * numFields

				query += "("
				for k := 0; k < numFields; k++ {
					query += `$` + strconv.Itoa(n+k+1) + ","
				}
				query = query[:len(query)-1] + "),"
				counter++
			}
		}
		query = query[:len(query)-1]
		_, err := db.Exec(query, values...)
		if err != nil {
			log.Println(query)
			return err
		}
	}
	return nil
}

// AToInt format time format to int
func AToInt(input string) int {
	result, _ := strconv.Atoi(input)
	return result
}
