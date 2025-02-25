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
	// Retrieves all sensors from the database.
	ListSensors(ctx context.Context, userID uuid.UUID, search string) ([]domain.Sensor, error)
	// Mark/uncheck sensors as favorites
	MarkSensorFavorite(ctx context.Context, userUuid uuid.UUID, sensorUuid uuid.UUID, favorite bool) error
}

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

func (r *SensorRepositoryImpl) ListSensors(ctx context.Context, userID uuid.UUID, search string) ([]domain.Sensor, error) {
	query := `
		SELECT uuid, name, category, description, visibility, SensorOwnerUuid
		FROM sensors
		WHERE (visibility = 1 OR SensorOwnerUuid = @userUuid) 
		AND (@search IS NULL OR name LIKE '%' + @search + '%')
	`

	rows, err := r.DB.QueryContext(ctx, query,
		sql.Named("userUuid", userID),
		sql.Named("search", search),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list sensors: %v", err)
	}
	defer rows.Close()

	var sensors []domain.Sensor
	for rows.Next() {
		var sensor domain.Sensor
		if err := rows.Scan(&sensor.ID, &sensor.Name, &sensor.Category, &sensor.Description, &sensor.Visibility, &sensor.SensorOwnerUuid); err != nil {
			return nil, fmt.Errorf("failed to scan sensor: %v", err)
		}
		sensors = append(sensors, sensor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sensors: %v", err)
	}

	return sensors, nil
}

func (r *SensorRepositoryImpl) MarkSensorFavorite(ctx context.Context, userUuid uuid.UUID, sensorUuid uuid.UUID, favorite bool) error {
	var query string

	if favorite {
		// Insert if the sensor is not already marked as a favorite for this user
		query = `
		IF NOT EXISTS (SELECT 1 FROM user_favorite_sensors WHERE userUuid = @userUuid AND sensorUuid = @sensorUuid)
		BEGIN
			INSERT INTO user_favorite_sensors (userUuid, sensorUuid)
			VALUES (@userUuid, @sensorUuid);
		END`
	} else {
		// Remove from the favorites if the user wants to unmark it
		query = `
		DELETE FROM user_favorite_sensors
		WHERE userUuid = @userUuid AND sensorUuid = @sensorUuid;
		`
	}

	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("userUuid", userUuid),
		sql.Named("sensorUuid", sensorUuid),
	)
	if err != nil {
		return fmt.Errorf("failed to update favorite status: %w", err)
	}

	return nil
}
