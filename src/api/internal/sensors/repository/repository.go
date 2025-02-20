package repository

import (
	config "api/configs"
	"api/internal/sensors/domain"
	"context"
	"fmt"

	"database/sql"
)

// Interface for sensor's data operations
type SensorRepository interface {
	// Creates a new sensor
	CreateSensor(ctx context.Context, sensor *domain.Sensor) error
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
