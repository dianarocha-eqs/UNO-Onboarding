package repository

import (
	config "api/configs"
	"api/internal/sensors/domain"
	"context"
	"fmt"

	"database/sql"

	_ "github.com/denisenkom/go-mssqldb" // Import SQL Server driver
)

// Interface for sensor's data operations
type SensorRepository interface {
	// Creates a new sensor
	CreateSensor(ctx context.Context, sensor *domain.Sensor) error
	// DeleteSensor removes a sensor from the database by its ID.
	DeleteSensor(id uint) error
	// GetAllSensors retrieves all sensors from the database.
	GetAllSensors() ([]domain.Sensor, error)
	// GetSensorByID retrieves a sensor by its ID from the database.
	GetSensorByID(id uint) (domain.Sensor, error)
	// UpdateSensor updates the details of an existing sensor in the database.
	UpdateSensor(sensor *domain.Sensor) error
}

// Performs user's data operations using database/sql to interact with the database
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

func (r *SensorRepositoryImpl) CreateSensor(ctx context.Context, sensor *domain.Sensor) error {
	query := `
		INSERT INTO Sensors (id, name, category, color, description, visibility, sensor_owner)
		VALUES (@id, @name, @category, @color, @description, @visibility, @sensor_owner)
	`

	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("id", sensor.ID),
		sql.Named("name", sensor.Name),
		sql.Named("category", sensor.Category),
		sql.Named("color", sensor.Color),
		sql.Named("description", sensor.Description),
		sql.Named("visibility", sensor.Visibility),
		sql.Named("sensor_owner", sensor.SensorOwner),
	)

	if err != nil {
		return fmt.Errorf("failed to create sensor: %v", err)
	}
	return nil
}

func (r *SensorRepositoryImpl) DeleteSensor(id uint) error {
	query := "DELETE FROM sensors WHERE id = ?"
	_, err := r.DB.Exec(query, id)
	return err
}

func (r *SensorRepositoryImpl) GetAllSensors() ([]domain.Sensor, error) {
	query := "SELECT id, name, category, color, description, visibility FROM sensors"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []domain.Sensor
	for rows.Next() {
		var sensor domain.Sensor
		if err := rows.Scan(&sensor.ID, &sensor.Name, &sensor.Category, &sensor.Color, &sensor.Description, &sensor.Visibility); err != nil {
			return nil, err
		}
		sensors = append(sensors, sensor)
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
