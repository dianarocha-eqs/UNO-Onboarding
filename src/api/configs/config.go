package configs

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// Holds the database and email configuration details.
type Config struct {
	// DB holds the database configuration details
	DB struct {
		// Host is the database server address
		Host string `json:"host"`
		// Port is the port number where the database is accessible
		Port int `json:"port"`
		// User is the username for authenticating with the database
		User string `json:"user"`
		// Password is the password for authenticating with the database
		Password string `json:"password"`
		// Name is the name of the database to connect to
		Name string `json:"name"`
	} `json:"db"`

	// Email holds the configuration for sending emails
	Email struct {
		// From is the email address that will appear as the sender
		From string `json:"from"`
		// Password is the password for the email account used for sending emails
		Password string `json:"password"`
		// Host is the email server address (SMTP server)
		Host string `json:"host"`
		// Port is the port number for connecting to the email server
		Port int `json:"port"`
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
