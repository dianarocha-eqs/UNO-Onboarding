package repository

import (
	config "api/configs"
	"api/internal/sensors/domain"
	"context"
	"fmt"

	"database/sql"

	_ "github.com/denisenkom/go-mssqldb" // Import SQL Server driver
)

// SensorRepository defines the methods required to interact with the sensor data storage.
// It provides basic CRUD operations for managing sensors.
type SensorRepository interface {
	// CreateSensor adds a new sensor to the database.
	CreateSensor(sensor *domain.Sensor) error
	// DeleteSensor removes a sensor from the database by its ID.
	DeleteSensor(id uint) error
	// Retrieves all sensors from the database.
	ListSensors(ctx context.Context, search string) ([]domain.Sensor, error)
	// GetSensorByID retrieves a sensor by its ID from the database.
	GetSensorByID(id uint) (domain.Sensor, error)
	// UpdateSensor updates the details of an existing sensor in the database.
	UpdateSensor(sensor *domain.Sensor) error
}

// SensorRepositoryImpl is the implementation of the SensorRepository interface.
// It uses GORM as the database ORM to interact with the database.
type SensorRepositoryImpl struct {
	DB *sql.DB
}

func NewSensorRepository() (SensorRepository, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &SensorRepositoryImpl{DB: db}, nil
}
func (r *SensorRepositoryImpl) CreateSensor(sensor *domain.Sensor) error {
	query := "INSERT INTO sensors (name, category, color, description, visibility VALUES (?, ?, ?, ?, ? GETDATE())"
	_, err := r.DB.Exec(query, sensor.Name, sensor.Category, sensor.Color, sensor.Description, sensor.Visibility)
	return err
}

func (r *SensorRepositoryImpl) DeleteSensor(id uint) error {
	query := "DELETE FROM sensors WHERE id = ?"
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *SensorRepositoryImpl) ListSensors(ctx context.Context, search string) ([]domain.Sensor, error) {
	query := `
		SELECT id, name, category, description
		FROM Sensor
		WHERE visibility = true AND name LIKE '%' + @search + @search + '%'
	`

	rows, err := r.DB.QueryContext(ctx, query,
		sql.Named("search", search))
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %v", err)
	}
	defer rows.Close()

	var sensors []domain.Sensor
	for rows.Next() {
		var sensor domain.Sensor
		if err := rows.Scan(&sensor.ID, &sensor.Name, &sensor.Category, &sensor.Description); err != nil {
			return nil, fmt.Errorf("failed to scan sensor: %v", err)
		}
		sensors = append(sensors, sensor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %v", err)
	}

	return sensors, nil
}

func (r *SensorRepositoryImpl) GetSensorByID(id uint) (domain.Sensor, error) {
	query := "SELECT id, name, category, color, description, visibility FROM sensors WHERE id = ?"
	var sensor domain.Sensor
	row := r.DB.QueryRow(query, id)
	if err := row.Scan(&sensor.ID, &sensor.Name, &sensor.Category, &sensor.Color, &sensor.Description, &sensor.Visibility); err != nil {
		return sensor, err
	}
	return sensor, nil
}

func (r *SensorRepositoryImpl) UpdateSensor(sensor *domain.Sensor) error {
	query := "UPDATE sensors SET name = ?, category = ?, color = ? , description = ?, visibility = GETDATE() WHERE id = ?"
	_, err := r.DB.Exec(query, sensor.Name, sensor.Category, sensor.Color, sensor.Description, sensor.Visibility)
	return err
}
