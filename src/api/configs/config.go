package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// DBConfig holds the database configuration details required for establishing a connection.
type DBConfig struct {
	Host     string `json:"host"`     // Host is the address of the database server.
	Port     int    `json:"port"`     // Port is the port number on which the database server is running.
	User     string `json:"user"`     // User is the username used to authenticate with the database.
	Password string `json:"password"` // Password is the user's password for database authentication.
	Name     string `json:"name"`     // Name is the name of the database to connect to.
}

// Config represents the overall application configuration, including database settings.
type Config struct {
	DB DBConfig `json:"db"` // DB contains the database configuration.
}

// ConfigFilePath is the relative path to the configuration JSON file.
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
