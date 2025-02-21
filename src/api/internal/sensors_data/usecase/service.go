package usecase

import (
	"api/internal/sensors_data/domain"
	"api/internal/sensors_data/repository"
	"context"
	"fmt"
)

type SensorDataService interface {
	AddSensorData(ctx context.Context, sensorData *domain.SensorData) error
}

type SensorDataServiceImpl struct {
	Repo repository.SensorDataRepository
}

func NewSensorDataService(repo repository.SensorDataRepository) SensorDataService {
	return &SensorDataServiceImpl{Repo: repo}
}

func (s *SensorDataServiceImpl) AddSensorData(ctx context.Context, sensorData *domain.SensorData) error {
	var err = s.Repo.AddSensorData(ctx, sensorData)
	if err != nil {
		return fmt.Errorf("failed to create sensor data")
	}
	return nil
}
