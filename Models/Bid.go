package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	
)

// Bid admin tabel model
type Bid struct {
	gorm.Model
	Time int64
	Amount float64
	AuctionID uint
	UserID uint
	Name string
}