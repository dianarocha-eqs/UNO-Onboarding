package repository

import (
	config "api/configs"
	"api/internal/sensors/domain"
	"fmt"

	"gorm.io/gorm"
)

type SensorRepository interface {
	CreateSensor(sensor *domain.Sensor) error
	DeleteSensor(id uint) error
	GetAllSensors() ([]domain.Sensor, error)
	GetSensorByID(id uint) (domain.Sensor, error)
	UpdateSensor(sensor *domain.Sensor) error
}

type SensorRepositoryImpl struct {
	DB *gorm.DB
}

func NewSensorRepository() (SensorRepository, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &SensorRepositoryImpl{DB: db}, nil
}
func (r *SensorRepositoryImpl) CreateSensor(sensor *domain.Sensor) error {
	return r.DB.Create(sensor).Error
}

func (r *SensorRepositoryImpl) DeleteSensor(id uint) error {
	return r.DB.Delete(&domain.Sensor{}, id).Error
}

func (r *SensorRepositoryImpl) GetAllSensors() ([]domain.Sensor, error) {
	var sensors []domain.Sensor
	err := r.DB.Find(&sensors).Error
	return sensors, err
}

func (r *SensorRepositoryImpl) GetSensorByID(id uint) (domain.Sensor, error) {
	var sensor domain.Sensor
	err := r.DB.First(&sensor, id).Error
	return sensor, err
}

func (r *SensorRepositoryImpl) UpdateSensor(sensor *domain.Sensor) error {
	return r.DB.Save(sensor).Error
}
