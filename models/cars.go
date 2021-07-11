package models

import (
	"gorm.io/gorm"
)

// Entity
type Car struct {
	gorm.Model
	Url       string
	Ip        string
	Country   string
	Province  string
	City      string
	Longitude float64
	Latitude  float64
	Rtts      []Rtt
}

// Interfaces
type CarDAO interface {
	Create(car Car) (Car, error)
	List() ([]Car, error)
}
