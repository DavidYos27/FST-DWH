package repositories

import (
	"database/sql"
	"log"
	"models"
	"strconv"
)

type postgresRecon struct {
	Conn *sql.DB
}

// NewRecon return implement of post repository interface
func NewRecon(Conn *sql.DB) ReconData {
	return &postgresRecon{
		Conn: Conn,
	}
}

const (
	entries = 500
)

// InsertToMatch insert to Match table
func (db *postgresRecon) InsertToMatch(tags []*models.Match) error {
	log.Println(len(tags))
	numFields := 7
	stmt := "INSERT INTO fact.report_match (provider_id, product_id, pz_id, bs_id,amount, status,date_id) values "
	for i := 0; i <= len(tags)/entries; i++ {
		query := stmt
		values := []interface{}{}
		counter := 0
		if i == len(tags)/entries {
			for j := i * entries; j < len(tags); j++ {
				values = append(values, AToInt(tags[j].ProviderID), AToInt(tags[j].ProductID), tags[j].PzID, tags[j].BsID, tags[j].Amount, tags[j].Status, tags[j].DateID)
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
				values = append(values, AToInt(tags[j].ProviderID), AToInt(tags[j].ProductID), tags[j].PzID, tags[j].BsID, tags[j].Amount, tags[j].Status, tags[j].DateID)
				n := counter * numFields

				query += "("
				for k := 0; k < numFields; k++ {
					query += `$` + strconv.Itoa(n+k+1) + ","
				}
				query = query[:len(query)-1] + "),"
				counter++
			}
		}
		if query != stmt {
			query = query[:len(query)-1]
			_, err := db.Conn.Exec(query, values...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// InsertToNotMatch Insert data to Not Match Table
func (db *postgresRecon) InsertToNotMatch(tags []*models.NotMatch) error {
	numFields := 7
	stmt := "INSERT INTO fact.report_not_match (provider_id, product_id,source_id, key_id,amount, status,date_id) values"
	for i := 0; i <= len(tags)/entries; i++ {
		query := stmt
		values := []interface{}{}
		counter := 0
		if i == len(tags)/entries {
			for j := i * entries; j < len(tags); j++ {
				values = append(values, AToInt(tags[j].ProviderID), AToInt(tags[j].ProductID), tags[j].Source, tags[j].KeyID, tags[j].Amount, tags[j].Status, tags[j].DateID)
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
				values = append(values, AToInt(tags[j].ProviderID), AToInt(tags[j].ProductID), tags[j].Source, tags[j].KeyID, tags[j].Amount, tags[j].Status, tags[j].DateID)
				n := counter * numFields

				query += "("
				for k := 0; k < numFields; k++ {
					query += `$` + strconv.Itoa(n+k+1) + ","
				}
				query = query[:len(query)-1] + "),"
				counter++
			}
		}
		if query != stmt {
			query = query[:len(query)-1]
			_, err := db.Conn.Exec(query, values...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// AToInt format time format to int
func AToInt(input string) int {
	result, _ := strconv.Atoi(input)
	return result
}
