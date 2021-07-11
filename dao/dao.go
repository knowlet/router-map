package dao

import (
	"github.com/knowlet/router-map/models"
	"gorm.io/gorm"
)

type DAO struct {
	Car models.CarDAO
}

func NewDAO(db *gorm.DB) *DAO {
	return &DAO{
		Car: CarDAO{DB: db},
	}
}
