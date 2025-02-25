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
	// Creates a new sensor
	CreateSensor(ctx context.Context, sensor *domain.Sensor) error
	// Updates the details of an existing sensor
	EditSensor(ctx context.Context, sensor *domain.Sensor) error
	// Returns true if sensorID has the same owner as userID
	GetSensorOwner(ctx context.Context, sensorUuid uuid.UUID, userID uuid.UUID) (bool, error)
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
		INSERT INTO Sensors (uuid, name, category, color, description, visibility, sensorOwnerUuid)
		VALUES (@uuid, @name, @category, @color, @description, @visibility, @sensorOwnerUuid)
	`

	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("uuid", sensor.ID),
		sql.Named("name", sensor.Name),
		sql.Named("category", sensor.Category),
		sql.Named("color", sensor.Color),
		sql.Named("description", sensor.Description),
		sql.Named("visibility", sensor.Visibility),
		sql.Named("sensorOwnerUuid", sensor.SensorOwnerUuid),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *SensorRepositoryImpl) EditSensor(ctx context.Context, sensor *domain.Sensor) error {

	query := `
		UPDATE sensors
		SET 
			name = COALESCE(NULLIF(@name, ''), name),
			category = COALESCE(NULLIF(@category, ''), category),
			color = COALESCE(NULLIF(@color, ''), color),
			description = @description,
			visibility = COALESCE(NULLIF(@visibility, ''), visibility)
		WHERE uuid = @uuid
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
		return err
	}
	return nil
}

func (r *SensorRepositoryImpl) GetSensorOwner(ctx context.Context, sensorUuid uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `
		SELECT sensorOwnerUuid
		FROM sensors
		WHERE uuid = @sensorUuid
	`

	var sensorOwnerUuid uuid.UUID
	err := r.DB.QueryRowContext(ctx, query, sql.Named("sensorUuid", sensorUuid)).Scan(&sensorOwnerUuid)
	if err != nil {
		return false, err
	}
	return true, nil
}
