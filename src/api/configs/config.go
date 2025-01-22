package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Config struct {
	DB DBConfig `json:"db"`
}

const ConfigFilePath = "../configs/config.json"

// LoadConfig reads the config file and unmarshals it into a Config object
func LoadConfig() (Config, error) {
	var config Config

	// Read the file
	data, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		return config, fmt.Errorf("could not read config file: %v", err)
	}

	// Parse the JSON
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("could not parse config file: %v", err)
	}

	return config, nil
}

// ConnectDB establishes a connection to the database using the loaded configuration
func ConnectDB() (*gorm.DB, error) {
	// Load the database configuration
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %v", err)
	}

	// URL-encode the password to handle special characters
	encodedPassword := url.QueryEscape(config.DB.Password)

	// Create the DSN (Data Source Name)
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		config.DB.User, encodedPassword, config.DB.Host, config.DB.Port, config.DB.Name)

	// Connect to the database
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}
