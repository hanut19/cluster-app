package db

/*
import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(dataSource string) error {
	var err error
	DB, err = sql.Open("postgres", dataSource)
	if err != nil {
		return err
	}
	return DB.Ping()
}
*/

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type Config struct {
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DBName     string `json:"db_name"`
	SSLMode    string `json:"ssl_mode"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func InitDB(configPath string) error {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.SSLMode,
	)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	return DB.Ping()
}

func SetDB(database *sql.DB) {
	DB = database
}
