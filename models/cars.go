package models

import (
	"gorm.io/gorm"
)

// Entity
type Car struct {
	gorm.Model `json:"-"`
	Url        string  `gorm:"uniqueIndex" json:"url" binding:"required"`
	Ip         string  `json:"ip"`
	User       string  `json:"user"`
	Pass       string  `json:"pass"`
	Country    string  `json:"country"`
	Province   string  `json:"province"`
	City       string  `json:"city"`
	Longitude  float64 `json:"lng"`
	Latitude   float64 `json:"lat"`
	Rtts       []Rtt   `json:"rtts"`
	Vendor     string  `json:"vendor"`
	Protocol   string  `json:"protocol"`
}

// Interfaces
type CarDAO interface {
	Create(car Car) (Car, error)
	List() ([]Car, error)
	Delete(id uint) error
}
