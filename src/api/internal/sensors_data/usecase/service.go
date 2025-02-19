package usecase

import (
	"api/internal/sensors_data/domain"
	"api/internal/sensors_data/repository"
	"context"
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
	return s.Repo.AddSensorData(ctx, sensorData)
}
