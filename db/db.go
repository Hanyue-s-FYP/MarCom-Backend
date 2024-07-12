package db

import (
	"database/sql"
	"fmt"

	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

func GetDB() *sql.DB {
	// opens database if not open before (check by pinging the database and see if there is any error)
	if db == nil {
		cfg := utils.GetConfig()
		dbOpen, err := sql.Open("sqlite3", cfg.DB_PATH)
        db = dbOpen
		if err != nil {
			panic(fmt.Sprintf("unable to open connection to sqlite3 database (%s): %v", cfg.DB_PATH, err))
		}
	}
	if err := db.Ping(); err != nil {
		cfg := utils.GetConfig()
		db, err = sql.Open("sqlite3", cfg.DB_PATH)
		if err != nil {
			panic(fmt.Sprintf("unable to open connection to sqlite3 database (%s): %v", cfg.DB_PATH, err))
		}
	}
	return db
}
