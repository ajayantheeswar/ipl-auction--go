package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User user tabel model
type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"unique;not null"`
	Password string
	AuthType string
	Auctions []Auction `gorm:"foreignkey:userid"`
}
