package utils

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// faster access than os.Getenv
type Config struct {
	PORT    string
	HOST    string
	DB_PATH string
}

func loadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		confKV := strings.SplitN(scanner.Text(), ",", 2)
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

func NewConfig(filename string) *Config {
	fmt.Printf("Loading Config from file `%s`\r\n", filename)
	if err := loadConfig(filename); err != nil {
		config = Config{
			HOST:    os.Getenv("HOST"),
			PORT:    os.Getenv("PORT"),
			DB_PATH: os.Getenv("DB_PATH"),
		}

	} else {
		// failed to load config smh, just fallback to default (since rn no involve critical API keys)
		config = Config{
			HOST:    "localhost",
			PORT:    "8080",
			DB_PATH: "marcom.db",
		}
	}

	return &config
}

func GetConfig() *Config {
	if config == (Config{}) {
		// programmer error, should call NewConfig to initialize first
		slog.Error("config is not initialized")
		panic("config is not initialized")
	}
	return &config
}
