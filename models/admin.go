package models

import (
	"errors"

	"github.com/ahyaghoubi/buses-live-location/db"
	"github.com/ahyaghoubi/buses-live-location/utils"
)

type Admin struct {
	ID       int64
	Name     string
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

var adminName = "Amir"
var adminEmail = "amir@amir.com"
var adminPassword = "12345678"

func (admin *Admin) Create() error {
	query := `
	INSERT INTO admins(name, email, password)
	VALUES (?,?,?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	admin.Password, err = utils.HashPassword(admin.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(admin.Name, admin.Email, admin.Password)
	if err != nil {
		return err
	}

	adminId, err := result.LastInsertId()

	admin.ID = adminId

	return err
}

func (admin *Admin) ValidateCredentials() error {
	query := "SELECT id, password FROM admins WHERE email = ?"
	row := db.DB.QueryRow(query, admin.Email)

	var retrievedPassword string

	err := row.Scan(&admin.ID, &retrievedPassword)
	if err != nil {
		return err
	}

	passwordIsValid := utils.ComparePassword(admin.Password, retrievedPassword)

	if !passwordIsValid {
		return errors.New("invalid credentials")
	}

	return nil
}

func (admin *Admin) UpdatePassword() error {
	query := `
		UPDATE bus
		SET PASSWORD = ?
		WHERE id = ?
	`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(admin.Password, admin.ID)

	return err
}

func CreateFirstAdminIfNotExist() error {
	query := "SELECT id FROM admins"

	rows, err := db.DB.Query(query)
	if err != nil {
		return err
	}

	defer rows.Close()

	var adminsid []int64

	for rows.Next() {
		rows.Scan(&adminsid)
	}

	if len(adminsid) == 0 {
		var admin Admin = Admin{
			Name:     adminName,
			Email:    adminEmail,
			Password: adminPassword,
		}

		admin.Create()
	}

	return nil
}
