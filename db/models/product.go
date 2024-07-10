package models

import "database/sql"

type Product struct {
	ID          int
	Name        sql.NullString
	Description sql.NullString
	Price       sql.NullFloat64
	Cost        sql.NullFloat64
	BusinessID  int
}
