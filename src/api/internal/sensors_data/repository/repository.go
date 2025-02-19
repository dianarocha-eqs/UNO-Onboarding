package repository

import (
	config "api/configs"
	"api/internal/sensors_data/domain"
	"context"
	"fmt"

	"database/sql"

	_ "github.com/denisenkom/go-mssqldb" // Import SQL Server driver
)

type SensorDataRepository interface {
	// Add sensor data
	AddSensorData(ctx context.Context, sensorData *domain.SensorData) error
}

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

func (r *SensorDataRepositoryImpl) AddSensorData(ctx context.Context, sensorData *domain.SensorData) error {

	query := `
		INSERT INTO Sensor_Data (id, sensorUuid, timestamp, value)
		VALUES (@id, @sensorUuid, @timestamp, @value)
	`

	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("id", sensorData.ID),
		sql.Named("sensorUuid", sensorData.Sensor),
		sql.Named("timestamp", sensorData.Timestamp),
		sql.Named("value", sensorData.Value),
	)

	if err != nil {
		return fmt.Errorf("failed to create sensor data: %v", err)
	}
	return nil

}
