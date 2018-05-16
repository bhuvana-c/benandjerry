package model

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type IceCreamStore interface {
	Create(iceCream *IceCream) error
	List() ([]*IceCream, error)
	Delete(name string) error
	Get(name string) (*IceCream, error)
	Update(iceCream *IceCream) error
}
type IceCream struct {
	Name                  string   `json:"name"`
	ImageClosed           string   `json:"imageClosed"`
	ImageOpen             string   `json:"imageOpen"`
	Description           string   `json:"description"`
	Story                 string   `json:"story"`
	SourcingValues        []string `json:"sourcingValues"`
	Ingredients           []string `json:"ingredients"`
	AllergyInfo           string   `json:"allergyInfo"`
	DietaryCertifications string   `json:"dietaryCertifications"`
	ProductID             string   `json:"productId"`
}

type NotFoundError struct {
}

func (error *NotFoundError) Error() string {
	return ""
}

type IceCreams struct {
	db *sql.DB
}

func NewIceCreamStore(db *sql.DB) (IceCreamStore, error) {
	return &IceCreams{db: db}, nil
}

func (iceCreamStore *IceCreams) Create(iceCream *IceCream) error {
	sqlQuery := "insert into icecream values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)"
	_, err := iceCreamStore.db.Exec(sqlQuery, iceCream.Name, iceCream.ImageClosed, iceCream.ImageOpen, iceCream.Description, iceCream.Story, pq.Array(iceCream.SourcingValues), pq.Array(iceCream.Ingredients), iceCream.AllergyInfo, iceCream.DietaryCertifications, iceCream.ProductID)
	if err != nil {
		return errors.Wrap(err, "error while creating a new info")
	}

	return err
}

func (iceCreamStore *IceCreams) Update(iceCream *IceCream) error {
	sqlQuery := "update icecream set image_closed=$1,image_open=$2,description=$3,story=$4,sourcing_values=$5,ingredients=$6,allergy_info=$7,dietary_certifications=$8,product_id=$9 where name=$10"
	_, err := iceCreamStore.db.Exec(sqlQuery, iceCream.ImageClosed, iceCream.ImageOpen, iceCream.Description, iceCream.Story, pq.Array(iceCream.SourcingValues), pq.Array(iceCream.Ingredients), iceCream.AllergyInfo, iceCream.DietaryCertifications, iceCream.ProductID, iceCream.Name)
	if err != nil {
		return errors.Wrap(err, "error while updating")
	}

	if err == sql.ErrNoRows {
		return &NotFoundError{}
	}
	return err
}

func (iceCreamStore *IceCreams) Get(name string) (*IceCream, error) {
	var iceCream = &IceCream{}
	err := iceCreamStore.db.QueryRow("select * from icecream where name=$1", name).Scan(&iceCream.Name, &iceCream.ImageClosed, &iceCream.ImageOpen, &iceCream.Description, &iceCream.Story, pq.Array(&iceCream.SourcingValues), pq.Array(&iceCream.Ingredients), &iceCream.AllergyInfo, &iceCream.DietaryCertifications, &iceCream.ProductID)

	if err == sql.ErrNoRows {
		return nil, &NotFoundError{}
	}
	if err != nil {
		return nil, err
	}
	return iceCream, nil
}
func (iceCreamStore *IceCreams) Delete(name string) error {
	_, err := iceCreamStore.db.Exec("delete from icecream where name=$1", name)
	return err
}

func (iceCreamStore *IceCreams) List() ([]*IceCream, error) {
	rows, err := iceCreamStore.db.Query("select * from icecream")
	if err != nil {
		return nil, errors.Wrap(err, "error while querying")
	}
	defer rows.Close()

	var iceCreams = make([]*IceCream, 0)
	for rows.Next() {
		var iceCream = &IceCream{}

		if err := rows.Scan(&iceCream.Name, &iceCream.ImageClosed, &iceCream.ImageOpen, &iceCream.Description, &iceCream.Story, pq.Array(&iceCream.SourcingValues), pq.Array(&iceCream.Ingredients), &iceCream.AllergyInfo, &iceCream.DietaryCertifications, &iceCream.ProductID); err != nil {
			return nil, errors.Wrap(err, "error while scanning")
		}
		iceCreams = append(iceCreams, iceCream)
	}

	return iceCreams, nil
}
