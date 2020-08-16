package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// AdminUser admin tabel model
type AdminUser struct {
	gorm.Model
	Name string
	Email string `gorm:"unique;not null"`
	Password string
	AuthType string
	Auctions []Auction 
}