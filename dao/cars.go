package dao

import (
	"github.com/knowlet/router-map/models"
	"gorm.io/gorm"
)

type CarDAO struct {
	DB *gorm.DB
}

func (dao CarDAO) Create(car models.Car) (models.Car, error) {
	err := dao.DB.Create(&car).Error
	if err != nil {
		return models.Car{}, err
	}
	return car, nil
}

func (dao CarDAO) Update(car *models.Car) (*models.Car, error) {
	err := dao.DB.Save(car).Error
	if err != nil {
		return nil, err
	}
	return car, nil
}

func (dao CarDAO) Delete(id uint) error {
	err := dao.DB.Delete(models.Car{Model: gorm.Model{ID: id}}).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao CarDAO) Get(id uint) (models.Car, error) {
	queryModel := models.Car{Model: gorm.Model{ID: id}}
	err := dao.DB.First(&queryModel).Error
	if err != nil {
		return models.Car{}, err
	}
	return queryModel, nil
}

func (dao CarDAO) List() ([]models.Car, error) {
	queryModel := []models.Car{}
	err := dao.DB.Find(&queryModel).Error
	if err != nil {
		return nil, err
	}
	return queryModel, nil
}
