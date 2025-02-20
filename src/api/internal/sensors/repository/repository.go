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
	// Updates the details of an existing sensor
	UpdateSensor(ctx context.Context, sensor *domain.Sensor) error
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

func (r *SensorRepositoryImpl) UpdateSensor(ctx context.Context, sensor *domain.Sensor) error {

	query := `
		UPDATE Sensor
		SET 
			name = COALESCE(NULLIF(@name, ''), name),
			category = COALESCE(NULLIF(@email, ''), category),
			color = COALESCE(NULLIF(@color, ''), color),
			description = @description,
			visibility = COALESCE(NULLIF(@visibility, ''), visibility)
		WHERE id = @id
	`

	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("id", sensor.ID),
		sql.Named("name", sensor.Name),
		sql.Named("category", sensor.Category),
		sql.Named("color", sensor.Color),
		sql.Named("description", sensor.Description),
		sql.Named("visibility", sensor.Visibility),
	)
	if err != nil {
		return fmt.Errorf("failed to update sensor: %v", err)
	}
	return nil
}
