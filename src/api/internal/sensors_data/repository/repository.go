package repository

import (
	config "api/configs"
	"api/internal/sensors_data/domain"
	"context"
	"fmt"
	"time"

	"database/sql"

	uuid "github.com/tentone/mssql-uuid"
)

// Interface for sensor's data operations
type SensorDataRepository interface {
	// Retrieves sensor data within a specific time interval.
	GetSensorData(ctx context.Context, sensorUuid uuid.UUID, from, to time.Time) ([]domain.SensorData, error)
}

// Performs sensors's data operations using database/sql to interact with the database
type SensorDataRepositoryImpl struct {
	DB *sql.DB
}

func NewSensorDataRepository() (SensorDataRepository, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &SensorDataRepositoryImpl{DB: db}, nil
}

func (s *SensorDataRepositoryImpl) GetSensorData(ctx context.Context, sensorUuid uuid.UUID, from, to time.Time) ([]domain.SensorData, error) {

	query := `
		SELECT timestamp, value
		FROM SensorData
		WHERE sensorUuid = @sensorUuid
		AND timestamp BETWEEN @from AND @to
		ORDER BY timestamp
	`

	rows, err := s.DB.QueryContext(ctx, query, sql.Named("sensorUuid", sensorUuid), sql.Named("from", from), sql.Named("to", to))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sensor data: %v", err)
	}
	defer rows.Close()

	var sensorData []domain.SensorData
	for rows.Next() {
		var data domain.SensorData
		if err := rows.Scan(&data.Timestamp, &data.Value); err != nil {
			return nil, fmt.Errorf("failed to scan sensor data: %v", err)
		}
		sensorData = append(sensorData, data)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over rows: %v", err)
	}

	return sensorData, nil
}
