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
	return dao.DB.Delete(&models.Car{Model: gorm.Model{ID: id}}).Error
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

func (dao CarDAO) GetProvinces() []string {
	var provinces []string
	dao.DB.Table("cars").Select("distinct province").Order("province").Scan(&provinces)
	return provinces
}

func (dao CarDAO) GetCars(province string) ([]models.Car, error) {
	var cars []models.Car
	err := dao.DB.Where("province = ?", province).Find(&cars).Error
	if err != nil {
		return nil, err
	}
	return cars, nil
}
