package configs

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// Holds the database configuration details required for establishing a connection.
type Config struct {
	DB struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"db"`

	Email struct {
		From     string `json:"from"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
	} `json:"email"`
}

// ConfigFilePath is the relative path to the configuration JSON file.
const ConfigFilePath = "../configs/config.json"

// LoadConfig reads the config file and unmarshals it into a Config object
func LoadConfig() (Config, error) {
	var config Config

	// Read the configuration file
	data, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		return config, fmt.Errorf("could not read config file at %s: %v", ConfigFilePath, err)
	}

	// Parse the JSON
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("could not parse config JSON: %v", err)
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
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Test the database connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}
	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database connection test failed: %v", err)
	}

	return db, nil
}
