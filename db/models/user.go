package models

import (
	"database/sql"
	"errors"

	"github.com/Hanyue-s-FYP/Marcom-Backend/db"
)

type User struct {
	ID          int
	Username    string
	Password    string `json:"-"` // always ignore this field in JSON encoding responses
	DisplayName string
	Email       string
	Status      int
	PhoneNumber string
}

type Investor struct {
	User
	ID int
}

type Business struct {
	User
	ID           int
	Description  string
	BusinessType string
	CoverImgPath string
}

type UserRole int

const (
	INVESTOR = iota
	BUSINESS
)

type LoginUserSqlResponse struct {
	User
	Role UserRole
}

// yes another workaround in go for without need to create another folder for another package namespace
// workaround so I can use BusinessModel.Create :)))
type businessModel struct{}

var BusinessModel *businessModel

func (*businessModel) Create(b Business) error {
	tx, err := db.GetDB().Begin()

	if err != nil {
		return err
	}

	userQuery := `
        INSERT INTO Users (username, password, display_name, email, status, phone_number)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	res, err := tx.Exec(userQuery, b.Username, b.Password, b.DisplayName, b.Email, b.Status, b.PhoneNumber)
	if err != nil {
		tx.Rollback()
		return err
	}

	userID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	businessQuery := `
			INSERT INTO Businesses (description, business_type, cover_img_path, user_id)
			VALUES (?, ?, ?, ?)
    `
	_, err = tx.Exec(businessQuery, b.Description, b.BusinessType, b.CoverImgPath, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

var ErrBusinessNotFound error = errors.New("business not found")

func (*businessModel) GetByBusinessID(id int) (*Business, error) {
	query := `
			SELECT
					u.id, u.username, u.password, u.display_name, u.email, u.status, u.phone_number,
					b.id, b.description, b.business_type, b.cover_img_path
			FROM
					Users u
			JOIN
					Businesses b ON u.id = b.user_id
			WHERE
					b.id = ?
    `
	row := db.GetDB().QueryRow(query, id)
	var business Business
	err := row.Scan(
		&business.User.ID, &business.Username, &business.Password, &business.DisplayName,
		&business.Email, &business.Status, &business.PhoneNumber,
		&business.ID, &business.Description, &business.BusinessType, &business.CoverImgPath,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBusinessNotFound
		}
		return nil, err
	}

	return &business, nil
}

func (*businessModel) GetByUsername(username string) (*Business, error) {
	query := `
			SELECT
					u.id, u.username, u.password, u.display_name, u.email, u.status, u.phone_number,
					b.id, b.description, b.business_type, b.cover_img_path
			FROM
					Users u
			JOIN
					Businesses b ON u.id = b.user_id
			WHERE
					u.username = ?
    `
	row := db.GetDB().QueryRow(query, username)
	var business Business
	err := row.Scan(
		&business.User.ID, &business.Username, &business.Password, &business.DisplayName,
		&business.Email, &business.Status, &business.PhoneNumber,
		&business.ID, &business.Description, &business.BusinessType, &business.CoverImgPath,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBusinessNotFound
		}
		return nil, err
	}

	return &business, nil
}

func (*businessModel) Update(business Business) error {
    // only description and cover image can be updated
    updateQuery := `
        UPDATE Businesses
        SET description = ?, cover_img_path = ?
        WHERE id = ?
    `

    _, err := db.GetDB().Exec(updateQuery, business.Description, business.CoverImgPath, business.ID)
    return err
}

type userModel struct{}

var UserModel *userModel

var ErrUserNotFound error = errors.New("user not found")

func (*userModel) GetByID(id int) (*User, error) {
	query := `
			SELECT
					u.id, u.username, u.password, u.display_name, u.email, u.status, u.phone_number
			FROM
					Users u
			WHERE
					u.id = ?
    `
	row := db.GetDB().QueryRow(query, id)
	var user User
	err := row.Scan(
		&user.ID, &user.Username, &user.Password, &user.DisplayName,
		&user.Email, &user.Status, &user.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}


