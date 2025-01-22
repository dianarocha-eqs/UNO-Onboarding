package usecase

import (
	"api/internal/sensors/domain"
	"api/internal/sensors/repository"
	"errors"
)

type SensorService interface {
	CreateSensor(sensor *domain.Sensor) error
	DeleteSensor(id uint) error
	GetAllSensors() ([]domain.Sensor, error)
	GetSensorByID(id uint) (domain.Sensor, error)
	UpdateSensor(sensor *domain.Sensor) error
}

type SensorServiceImpl struct {
	Repo repository.SensorRepository
}

func NewSensorService(repo repository.SensorRepository) SensorService {
	return &SensorServiceImpl{Repo: repo}
}

func (s *SensorServiceImpl) CreateSensor(sensor *domain.Sensor) error {
	if sensor.Name == "" && sensor.Category == "" {
		return errors.New("name and category are required")
	}
	return s.Repo.CreateSensor(sensor)
}

func (s *SensorServiceImpl) UpdateSensor(sensor *domain.Sensor) error {
	if sensor.Name == "" && sensor.Category == "" {
		return errors.New("name and category are required")
	}
	return s.Repo.UpdateSensor(sensor)
}

func (s *SensorServiceImpl) DeleteSensor(id uint) error {
	return s.Repo.DeleteSensor(id)
}

func (s *SensorServiceImpl) GetAllSensors() ([]domain.Sensor, error) {
	return s.Repo.GetAllSensors()
}

func (s *SensorServiceImpl) GetSensorByID(id uint) (domain.Sensor, error) {
	return s.Repo.GetSensorByID(id)
}
