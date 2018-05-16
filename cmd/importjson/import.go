package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type IceCream struct {
	Name                  string   `json:"name"`
	ImageClosed           string   `json:"image_closed"`
	ImageOpen             string   `json:"image_open"`
	Description           string   `json:"description"`
	Story                 string   `json:"story"`
	SourcingValues        []string `json:"sourcing_values"`
	Ingredients           []string `json:"ingredients"`
	AllergyInfo           string   `json:"allergy_info"`
	DietaryCertifications string   `json:"dietary_certifications"`
	ProductID             string   `json:"productId"`
}

func initDB(dbURL string, jsonFile string) error {

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		errors.Wrap(err, "error while connecting to database")
	}
	defer db.Close()

	//	absPath, err := filepath.Abs("../benandjerry/icecream.json")
	//	if err != nil {
	//		return errors.Wrap(err, "error while finding absoulte path")
	//	}
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return errors.Wrap(err, "error while reading icecream.json")
	}
	t := time.Now()
	var iceCreams = make([]*IceCream, 0)
	err = json.Unmarshal(data, &iceCreams)
	if err != nil {
		return errors.Wrap(err, "error while decoding json data")
	}
	fmt.Println("======", len(iceCreams))
	fmt.Println(time.Since(t).Seconds)
	txn, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "error while creating a transaction")
	}

	stmt, err := txn.Prepare(pq.CopyIn("icecream", "name", "image_closed", "image_open", "description", "story", "sourcing_values", "ingredients", "allergy_info", "dietary_certifications", "product_id"))
	if err != nil {
		return errors.Wrap(err, "error while preparing a transaction")
	}

	for _, iceCream := range iceCreams {
		fmt.Println("sss", pq.Array(iceCream.Ingredients), iceCream.ProductID)
		_, err = stmt.Exec(iceCream.Name, iceCream.ImageClosed, iceCream.ImageOpen, iceCream.Description, iceCream.Story, pq.Array(iceCream.SourcingValues), pq.Array(iceCream.Ingredients), iceCream.AllergyInfo, iceCream.DietaryCertifications, iceCream.ProductID)
		if err != nil {
			return errors.Wrap(err, "error while adding items fo bulk import")
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return errors.Wrap(err, "error while performing bulk import")
	}

	err = stmt.Close()
	if err != nil {
		return errors.Wrap(err, "error while closing bulk import")
	}

	err = txn.Commit()
	if err != nil {
		return errors.Wrap(err, "error while closing transaction")
	}
	fmt.Println("done")
	return nil

}

func main() {

	connStrFormat := "postgres://{username}:{password}@{hostname}:{port}/{dbname}"
	dbURL := flag.String("db", "", "Database connection string with the format - "+connStrFormat)
	jsonFile := flag.String("file", "", "Data to be imported")

	flag.Parse()

	if dbURL == nil || *dbURL == "" {
		log.Fatalf("argument -db is missing")
	}
	if jsonFile == nil || *jsonFile == "" {
		log.Fatalf("argument -file is missing")
	}
	fmt.Println(*dbURL, *jsonFile)
	err := initDB(*dbURL, *jsonFile)
	if err != nil {
		log.Fatalf("error while importing data: %s", err)
	}
}
