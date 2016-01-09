package main

import (
	"database/sql"
	"fmt"
	"github.com/benmanns/goworker"
	"github.com/kr/pretty"
	"log"
	"time"
	// "os"
)

import _ "github.com/lib/pq"

type curationTask struct {
	ID                  int
	CustomerID          int
	CurationJobID       int
	StrategyAttributeID int
	ProductID           int
	Status              string
	CreatedAt           time.Time
}

func unpackInterfaceMap(args map[string]interface{}) map[string]string {
	result := map[string]string{}

	for key, value := range args {
		result[key] = value.(string)
	}
	return result
}

func curateColor(queue string, arguments ...interface{}) error {
	args := unpackInterfaceMap(arguments[0].(map[string]interface{}))

	fmt.Println("args:")
	pretty.Print(args)
	fmt.Printf("\n\n")

	// db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	db, err := sql.Open("postgres", "host=192.168.59.103 user=postgres password=postgres dbname=feeddata sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`SELECT id, customer_id, curation_job_id, strategy_attribute_id, product_id, status, created_at from curation_tasks WHERE curation_job_id=$1`, args["task_id"])
	if err != nil {
		log.Fatal(err)
	}

	results := map[int]curationTask{}
	for rows.Next() {
		//temp vars
		var tempId int
		var tempCustomerId int
		var tempCurationJobId int
		var tempStrategyAttributeId int
		var tempProductId int
		var tempStatus string
		var tempCreatedAt time.Time

		if err := rows.Scan(&tempId, &tempCustomerId, &tempCurationJobId, &tempStrategyAttributeId, &tempProductId, &tempStatus, &tempCreatedAt); err != nil {
			log.Fatal(err)
		}
		dbResult := curationTask{
			ID:                  tempId,
			CustomerID:          tempCustomerId,
			CurationJobID:       tempCurationJobId,
			StrategyAttributeID: tempStrategyAttributeId,
			ProductID:           tempProductId,
			Status:              tempStatus,
			CreatedAt:           tempCreatedAt,
		}
		results[tempId] = dbResult
		fmt.Println(dbResult)
	}

	// pretty.Print(results)

	return nil
}

func init() {
	goworker.Register("ColorAnalysis", curateColor)
}

func main() {
	if err := goworker.Work(); err != nil {
		fmt.Println("Error:", err)
	}
}
