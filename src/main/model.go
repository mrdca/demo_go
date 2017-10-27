package main

import (
	"database/sql"
)

type smartphone struct{
	ID int `json:"id"`
	Name string `json:"name"`
	Brand string `json:"brand"`
	Price float64 `json:"price"`
}

func (s *smartphone) getSmartphone(db *sql.DB) error {
	return db.QueryRow("SELECT name, price, brand FROM smartphone WHERE id=$1",
		s.ID).Scan(&s.Name, &s.Price, &s.Brand)
}

func (s *smartphone) updateSmartphone(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE smartphone SET name=$1, price=$2, brand=$3 WHERE id=$4",
			s.Name, s.Price, s.Brand, s.ID)

	return err
}

func (s *smartphone) deleteSmartphone(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM smartphone WHERE id=$1", s.ID)

	return err
}

func (s *smartphone) createSmartphone(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO smartphone(name, brand, price ) VALUES($1, $2, $3) RETURNING id",
		s.Name, s.Brand, s.Price).Scan(&s.ID)

	if err != nil {
		return err
	}

	return nil
}

func getSmartphones(db *sql.DB, start, count int) ([]smartphone, error) {
	rows, err := db.Query(
		"SELECT id, name, price, brand FROM smartphone LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	smartphones := []smartphone{}

	for rows.Next() {
		var s smartphone
		if err := rows.Scan(&s.ID, &s.Name, &s.Price, &s.Brand); err != nil {
			return nil, err
		}
		smartphones = append(smartphones, s)
	}

	return smartphones, nil
}
