package db

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/Hanyue-s-FYP/Marcom-Backend/utils"
	"github.com/google/uuid"
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

// return the final path (would be the same as in physical path and served file)
func UploadImage(file *multipart.File, header *multipart.FileHeader) (string, error) {
	// get where to store files
	path := fmt.Sprintf("%s/", utils.GetConfig().IMG_FOLDER)
	slog.Debug(fmt.Sprintf("Path configured: %s", utils.GetConfig().IMG_FOLDER))

	id := uuid.New()
	fileExtension := strings.Split(header.Filename, ".")
	if len(fileExtension) < 2 {
		return "", utils.HttpError{
			Code:       http.StatusBadRequest,
			Message:    "Please provide a file with image extension (eg. png)",
			LogMessage: "file does not have extension",
		}
	}

	// possible file has more than 1 dots (not possible for avg user though), so will just take the last extension (no checks for now ;))
	uploadPath := path + id.String() + "." + fileExtension[len(fileExtension)-1]
	f, err := os.OpenFile(uploadPath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", err
	}

	// copy the file to the physical path on disk
	io.Copy(f, *file)

	return uploadPath, nil
}
