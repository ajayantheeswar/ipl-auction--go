package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"ipl/models"
	"fmt"
)

// Db Database
var Db *gorm.DB

func InitialiseDatabase () {
	db, err := gorm.Open("postgres", "host port=5432 user= dbname= password=")
  	if err != nil {
		panic(err)
	}else{
		Db = db
	}

	/*db.Debug().DropTableIfExists(&models.AdminUser{})
	db.Debug().DropTableIfExists(&models.User{})
	db.Debug().DropTableIfExists(&models.Auction{})
	db.Debug().DropTableIfExists(&models.Bid{})*/

	db.Debug().AutoMigrate(&models.AdminUser{})
	db.Debug().AutoMigrate(&models.User{})
	db.Debug().AutoMigrate(&models.Auction{})
	db.Debug().AutoMigrate(&models.Bid{})

}

func BidWatchDogFucntion() {
	Db.Exec("SELECT stopbid();");
	fmt.Print("Executed");
}
