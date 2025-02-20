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
	// Retrieves all sensors from the database.
	ListSensors(ctx context.Context, userID uuid.UUID, search string) ([]domain.Sensor, error)
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

func (r *SensorRepositoryImpl) ListSensors(ctx context.Context, userID uuid.UUID, search string) ([]domain.Sensor, error) {
	query := `
		SELECT id, name, category, description, visibility, sensor_owner
		FROM sensors
		WHERE (visibility = 0 OR sensor_owner = @user_id) 
		  AND (@search = '' OR name LIKE '%' + @search + '%')
	`

	rows, err := r.DB.QueryContext(ctx, query,
		sql.Named("user_id", userID),
		sql.Named("search", search),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list sensors: %v", err)
	}
	defer rows.Close()

	var sensors []domain.Sensor
	for rows.Next() {
		var sensor domain.Sensor
		if err := rows.Scan(&sensor.ID, &sensor.Name, &sensor.Category, &sensor.Description, &sensor.Visibility, &sensor.SensorOwner); err != nil {
			return nil, fmt.Errorf("failed to scan sensor: %v", err)
		}
		sensors = append(sensors, sensor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sensors: %v", err)
	}

	return sensors, nil
}
