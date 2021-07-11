package models

import (
	"gorm.io/gorm"
)

type Rtt struct {
	gorm.Model
	Min    float64
	Avg    float64
	Max    float64
	Online bool
	CarID  uint
}
