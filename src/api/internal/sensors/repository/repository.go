package repository

import (
	config "api/configs"
	"context"
	"fmt"

	"database/sql"

	uuid "github.com/tentone/mssql-uuid"
)

type SensorRepository interface {
	MarkSensorFavorite(ctx context.Context, sensorUuid uuid.UUID, favorite bool) error
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

func (s *SensorRepositoryImpl) MarkSensorFavorite(ctx context.Context, sensorUuid uuid.UUID, favorite bool) error {
	query := `
        UPDATE Sensors
        SET favorite = @favorite
        WHERE id = @sensorUuid
    `

	_, err := s.DB.ExecContext(ctx, query, sql.Named("favorite", favorite), sql.Named("sensorUuid", sensorUuid))
	if err != nil {
		return fmt.Errorf("failed to update favorite status: %w", err)
	}

	return nil
}
