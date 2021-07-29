package models

import (
	"gorm.io/gorm"
)

// Entity
type Car struct {
	gorm.Model
	Url       string  `json:"url" binding:"required"`
	Ip        string  `json:"ip"`
	User      string  `json:"user" binding:"required"`
	Pass      string  `json:"pass" binding:"required"`
	Country   string  `json:"country"`
	Province  string  `json:"province"`
	City      string  `json:"city"`
	Longitude float64 `json:"lng"`
	Latitude  float64 `json:"lat"`
	Rtts      []Rtt   `json:"rtts"`
}

// Interfaces
type CarDAO interface {
	Create(car Car) (Car, error)
	List() ([]Car, error)
}
