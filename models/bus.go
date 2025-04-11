package models

import (
	"github.com/ahyaghoubi/buses-live-location/db"
)

type Bus struct {
	ID   int64
	Name string `binding:"required"`
}

func (b *Bus) Create() error {
	query := `
		INSERT INTO bus (name)
		VALUES (?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(b.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	b.ID = id
	return err
}

func (b *Bus) Update() error {
	query := `
		UPDATE bus
		SET name = ?
		WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(b.Name, b.ID)
	return err
}

func (b *Bus) Delete() error {
	query := "DELETE FROM bus WHERE id = ?"

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(b.ID)
	return err
}

func GetAllBuses() ([]Bus, error) {
	query := "SELECT * FROM bus"

	rows, err := db.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var buses []Bus

	for rows.Next() {
		var bus Bus
		err := rows.Scan(&bus.ID, &bus.Name)
		if err != nil {
			return nil, err
		}

		buses = append(buses, bus)
	}

	return buses, nil
}

func GetBusById(id int64) (*Bus, error) {
	query := "SELECT * FROM bus WHERE id = ?"

	row := db.DB.QueryRow(query, id)

	var bus Bus

	err := row.Scan(&bus.ID, &bus.Name)
	if err != nil {
		return nil, err
	}

	return &bus, nil
}
