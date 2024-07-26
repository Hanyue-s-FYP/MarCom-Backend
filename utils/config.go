package utils

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// faster access than os.Getenv
type Config struct {
	PORT       string
	HOST       string
	DB_PATH    string
	IMG_FOLDER string
}

func loadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
        confKV := strings.SplitN(scanner.Text(), "=", 2)
		if len(confKV) != 2 {
			continue
		}
		key, value := confKV[0], confKV[1]
		os.Setenv(key, value)
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}

// store as global so when getting config for the second time no need to create a new one
var config Config
func (config *Config) String() string {
    return fmt.Sprintf("HOST: %s, PORT: %s, DB_PATH: %s, IMG_PATH: %s", config.HOST, config.PORT, config.DB_PATH, config.IMG_FOLDER)
}

func NewConfig(filename string) *Config {
	fmt.Printf("Loading Config from file `%s`\r\n", filename)
	if err := loadConfig(filename); err == nil {
		config = Config{
			HOST:       os.Getenv("HOST"),
			PORT:       os.Getenv("PORT"),
			DB_PATH:    os.Getenv("DB_PATH"),
			IMG_FOLDER: os.Getenv("IMG_FOLDER"),
		}
	} else {
        slog.Error(fmt.Sprintf("failed to load config, using default values: %v", err))
		// failed to load config smh, just fallback to default (since rn no involve critical API keys)
		config = Config{
			HOST:       "localhost",
			PORT:       "8080",
			DB_PATH:    "marcom.db",
			IMG_FOLDER: "img",
		}
	}
    slog.Info(fmt.Sprintf("Config loaded: %s", config.String()))

    // create img dir if not exists
    imgPath := filepath.Join(".", config.IMG_FOLDER)
    if err := os.MkdirAll(imgPath, os.ModePerm); err != nil {
        slog.Error(fmt.Sprintf("failed to check and create IMG_FOLDER: %v", err))
    }

	return &config
}

func GetConfig() *Config {
	if config == (Config{}) {
		// programmer error, should call NewConfig to initialize first
        slog.Error(fmt.Sprintf("config is not initialized: %s", config.String()))
		panic("config is not initialized")
	}
	return &config
}
