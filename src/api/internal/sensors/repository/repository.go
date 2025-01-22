package repository

import (
	config "api/configs"
	"api/internal/sensors/domain"
	"fmt"

	"gorm.io/gorm"
)

// SensorRepository defines the methods required to interact with the sensor data storage.
// It provides basic CRUD operations for managing sensors.
type SensorRepository interface {
	// CreateSensor adds a new sensor to the database.
	CreateSensor(sensor *domain.Sensor) error
	// DeleteSensor removes a sensor from the database by its ID.
	DeleteSensor(id uint) error
	// GetAllSensors retrieves all sensors from the database.
	GetAllSensors() ([]domain.Sensor, error)
	// GetSensorByID retrieves a sensor by its ID from the database.
	GetSensorByID(id uint) (domain.Sensor, error)
	// UpdateSensor updates the details of an existing sensor in the database.
	UpdateSensor(sensor *domain.Sensor) error
}

// SensorRepositoryImpl is the implementation of the SensorRepository interface.
// It uses GORM as the database ORM to interact with the database.
type SensorRepositoryImpl struct {
	DB *gorm.DB // DB is the GORM instance used to interact with the database.
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
