package dao

import (
	"fmt"
	"strings"

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

func Migration(db *gorm.DB) {
	// create extension
	// createExtension(db, "uuid-ossp")

	// create enmu type

	// auto migration
	db.AutoMigrate(&models.Car{})

}

func createExtension(db *gorm.DB, ext string) {
	db.Exec(fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS \"%s\";", ext))
}

func createEnum(db *gorm.DB, enumname string, enumvalues []string) {
	db.Exec(fmt.Sprintf(
		"CREATE TYPE %s AS ENUM (%s);",
		enumname,
		"'"+strings.Join(enumvalues, "', '")+"'",
	))
}
