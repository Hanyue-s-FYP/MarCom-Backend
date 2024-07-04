package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// faster access than os.Getenv
type Config struct {
    PORT string
    HOST string
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

func NewConfig(filename string) *Config {
    var config Config

    fmt.Printf("Loading Config from file `%s`\r\n", filename)

    return &config
}
