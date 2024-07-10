package models

import "database/sql"

type Environment struct {
	Products    []Product
	Agents      []Agent
	ID          int
	Name        sql.NullString
	Description sql.NullString
	BusinessID  int
}
