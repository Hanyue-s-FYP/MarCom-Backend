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
	PORT           string
	HOST           string
	DB_PATH        string
	GRPC_CORE_ADDR string
	IMG_FOLDER     string
	JWT_SECRET_KEY string
	FRONT_END_ADDR string
	SMTP_ADDR      string
	SMTP_PORT      string
	SMTP_EMAIL     string
	SMTP_PW        string
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
    // only log those not secret stuff
	return fmt.Sprintf(
		"HOST: %s, PORT: %s, DB_PATH: %s, IMG_PATH: %s, GRPC_CORE_ADDR: %s, FRONT_END_ADDR: %s, SMTP_ADDR: %s, SMTP_PORT: %s",
		config.HOST,
		config.PORT,
		config.DB_PATH,
		config.IMG_FOLDER,
		config.GRPC_CORE_ADDR,
		config.FRONT_END_ADDR,
		config.SMTP_ADDR,
		config.SMTP_PORT,
	)
}

// if setup environment variables in system alrd, pass in empty string and will attempt to load environment variables from the system
func NewConfig(filename string) *Config {
	fmt.Printf("Loading Config from file `%s`\r\n", filename)

	if filename != "" {
		err := loadConfig(filename)
		if err != nil {
			panic("failed to load environment variables from " + filename)
		}
	}

	config = Config{
		HOST:           os.Getenv("HOST"),
		PORT:           os.Getenv("PORT"),
		DB_PATH:        os.Getenv("DB_PATH"),
		GRPC_CORE_ADDR: os.Getenv("GRPC_CORE_ADDR"),
		IMG_FOLDER:     os.Getenv("IMG_FOLDER"),
		JWT_SECRET_KEY: os.Getenv("JWT_SECRET_KEY"),
		FRONT_END_ADDR: os.Getenv("FRONT_END_ADDR"),
		SMTP_ADDR:      os.Getenv("SMTP_ADDR"),
		SMTP_PORT:      os.Getenv("SMTP_PORT"),
		SMTP_EMAIL:     os.Getenv("SMTP_EMAIL"),
		SMTP_PW:        os.Getenv("SMTP_PW"),
	}
	slog.Info(fmt.Sprintf("Config loaded: %s", config.String()))

	// if cant load SMTP stuffs, panic
	if config.SMTP_ADDR == "" || config.SMTP_PORT == "" || config.SMTP_EMAIL == "" || config.SMTP_PW == "" {
		slog.Error("SMTP not configured")
		panic("SMTP not configured")
	}

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
