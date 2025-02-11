package configs

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"database/sql"
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
	// JWTSecret holds the JWT secret key
	JWTSecret string `json:"jwt_secret"`
}

// ConfigFilePath is the relative path to the configuration JSON file.
const ConfigFilePath = "../configs/config.json"

// LoadConfig reads and parses the configuration file
func LoadConfig() (Config, error) {
	var config Config

	// Read the configuration file
	data, err := os.ReadFile(ConfigFilePath)
	if err != nil {
		return config, fmt.Errorf("could not read config file: %v", err)
	}

	// Parse JSON
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("could not parse config JSON: %v", err)
	}

	return config, nil
}

// ConnectDB establishes a connection to the SQL Server database
func ConnectDB() (*sql.DB, error) {
	// Load database configuration
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database config: %v", err)
	}

	// URL-encode password
	encodedPassword := url.QueryEscape(config.DB.Password)

	// Construct DSN (Data Source Name) with Azure-compliant parameters
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=true&TrustServerCertificate=false",
		config.DB.User, encodedPassword, config.DB.Host, config.DB.Port, config.DB.Name)

	// Connect to the database
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}
