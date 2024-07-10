package models

import "database/sql"

type User struct {
	ID          int
	Username    string
	Password    string
	DisplayName sql.NullString
	Email       sql.NullString
	Status      sql.NullString
	PhoneNumber sql.NullString
}

type Investor struct {
	User
	ID int
}

type Business struct {
	User
	ID           int
	Description  sql.NullString
	BusinessType sql.NullString
	CoverImgPath sql.NullString
}
