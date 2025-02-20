package repository

import (
	config "api/configs"
	"api/internal/sensors/domain"
	"context"
	"fmt"

	"database/sql"

	uuid "github.com/tentone/mssql-uuid"
)

// Interface for sensor's data operations
type SensorRepository interface {
	// Updates the details of an existing sensor
	UpdateSensor(ctx context.Context, sensor *domain.Sensor) error
	// Returns true if sensorID has the same owner as userID
	GetSensorOwner(ctx context.Context, sensorID uuid.UUID, userID uuid.UUID) (bool, error)
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
		UPDATE sensors
		SET 
			name = COALESCE(NULLIF(@name, ''), name),
			category = COALESCE(NULLIF(@category, ''), category),
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

func (r *SensorRepositoryImpl) GetSensorOwner(ctx context.Context, sensorID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `
		SELECT sensor_owner
		FROM sensors
		WHERE id = @sensorID
	`

	var sensorOwner uuid.UUID
	err := r.DB.QueryRowContext(ctx, query, sql.Named("sensorID", sensorID)).Scan(&sensorOwner)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("sensor not found")
		}
		return false, fmt.Errorf("failed to retrieve sensor owner: %v", err)
	}

	return sensorOwner == userID, nil
}
