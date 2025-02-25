package repository

import (
	config "api/configs"
	"api/internal/sensors_data/domain"
	"context"
	"fmt"

	"database/sql"
)

// Interface for sensor's data operations
type SensorDataRepository interface {
	// Add sensor data
	AddSensorData(ctx context.Context, sensorData []*domain.SensorData) error
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

func (r *SensorDataRepositoryImpl) AddSensorData(ctx context.Context, sensorData []*domain.SensorData) error {

	// Start a transaction
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	query := `
		INSERT INTO SensorData (sensorUuid, timestamp, value)
		VALUES (@sensorUuid, @timestamp, @value)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, sensorData := range sensorData {
		_, err := stmt.ExecContext(ctx,
			sql.Named("sensorUuid", sensorData.SensorUuid),
			sql.Named("timestamp", sensorData.Timestamp),
			sql.Named("value", sensorData.Value),
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert sensor data: %v", err)
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil

}
