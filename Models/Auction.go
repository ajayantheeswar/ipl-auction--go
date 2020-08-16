package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Auction user tabel model
type Auction struct {
	gorm.Model
	Name         string
	BattingStyle string 
	Average      string
	Role         string
	Start        uint64
	End          uint64
	IsStarted    bool `gorm:"default:false"`
	IsActive     bool `gorm:"default:false"`
	IsSold       bool `gorm:"default:false"`
	Profile      string
	Bids [] Bid
	AdminUserID uint
	UserID uint
}
