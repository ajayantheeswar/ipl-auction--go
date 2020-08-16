package controllers

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"ipl/models"
	"ipl/database"
	"ipl/utils"
	"github.com/jinzhu/gorm"
	"time"

)

func SignUpUser(c *gin.Context) {
	var requestBody POJOSignUP
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Auth Status" : "FAIL"})
		panic(err)
		return
	}

	var user = models.User {}

	if (requestBody.AuthType == "Google"){
		Payload,err := utils.VerifyIDToken(requestBody.Token)
		
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Auth Status": "FAIL"})
			panic(err)
			return
		}
		user.Email = Payload.Claims["email"].(string)
		user.Name = Payload.Claims["name"].(string)
		user.AuthType = "Google"
	}else{
		encodedPassword, err := utils.HashPassword(requestBody.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Auth Status": "FAIL"})
			panic(err)
			return
		}
		user.Email = requestBody.Email
		user.Name = requestBody.Name
		user.AuthType = "MH"
		user.Password = encodedPassword
	}

	// Check if the user Already Present
	err = database.Db.Create(&user).Error	

	if err != nil {
		// User EXISTS ... RETURN
		c.JSON(http.StatusConflict, gin.H{"Auth Status" : "FAIL" ,"error" : "The Account Already Exists"})
	}else{
		// User Not Exists
		token, err:= utils.CreateToken(fmt.Sprint(user.ID))
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"Auth Status" : "FAIL"})
			panic(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"Token" : token,
            "Name" : user.Name , "Email" : user.Email,
        })
	}

}


func SignInUser(c *gin.Context) {
	var requestBody POJOSignIN
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		panic(err)
		c.JSON(http.StatusBadRequest, gin.H{"Auth Status": "FAIL"})
		return
	}
	var user models.User

	if (requestBody.AuthType == "Google"){
		Payload,err := utils.VerifyIDToken(requestBody.Token)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"Auth Status": "FAIL"})
			panic(err)
			return
		}
		requestBody.Email = Payload.Claims["email"].(string)
	}

	// Check if the user Already Present
	err = database.Db.Where(&models.User{Email: requestBody.Email}).First(&user).Error

	if gorm.IsRecordNotFoundError(err) {
		// User Does Not EXISTS ... RETURN
		c.JSON(http.StatusForbidden, gin.H{"Auth Status": "FAIL", "error": "Invalid Credientials"})
		panic(err)
		return
	} else {
		// User Exists
		if(requestBody.AuthType != "Google"){
			if !utils.CheckPasswordHash(requestBody.Password, user.Password) {
				// Invalid Credientials
				c.JSON(http.StatusForbidden, gin.H{"Auth Status": "FAIL", "error": "Invalid Credientials"})
				return
			}
		}
		token, err := utils.CreateToken(fmt.Sprint(user.ID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Auth Status": "FAIL"})
			panic(err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"Token": token,
			"Name": user.Name, "Email": user.Email,
		})
	}
}

type POJOCreateBid struct {
    Amount float64
    AuctionID uint
}

func UserCreateBid(c *gin.Context) {

	var requestBody POJOCreateBid
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Auth Status" : "FAIL"})
		return
	}
	user := c.Request.Context().Value("UserId").(models.User)
	Bid := models.Bid {
			Amount : requestBody.Amount,
			AuctionID : requestBody.AuctionID,
			UserID : user.ID,
			Name : user.Name,
		} 
	
	Bid.Time = getTime()
	err = database.Db.Create(&Bid).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status" : "Creation Failed"})
		return
	}

	Auction := models.Auction{}
	err = database.Db.Debug().Preload("Bids").First(&Auction,requestBody.AuctionID).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status" : "Fetch Failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status" : "Createion SuccessFul","Auction": Auction});
}

func getTime() (int64){
	now := time.Now()
	unixNano := now.UnixNano()                                                                      
	umillisec := unixNano / 1000000 
	return umillisec
}

func UserGetAllAuctions(c *gin.Context) {
	
	var Auctions [] models.Auction

	err := database.Db.Debug().Where(models.Auction{IsSold : false}).Find(&Auctions).Error
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Fetch FAILED"})	
		return
	}
	c.JSON(http.StatusOK, gin.H{"Auctions": Auctions});

}

type POJOUserGetAuction struct {
	AuctionID uint
}

func UserGetAuction(c *gin.Context) {
	user := c.Request.Context().Value("UserId").(models.User)

	var requestBody POJOAdminGetAuction
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Request FAILED"})	
	}

	var Auction models.Auction
	err = database.Db.Debug().Preload("Bids","user_id = ?",user.ID).First(&Auction,requestBody.AuctionID).Error
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Fetch FAILED"})	
		panic(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"Auction": Auction});

}

func UserGetAllUserAuction (c *gin.Context) {
	user := c.Request.Context().Value("UserId").(models.User)

	var Auctions []models.Auction
	err := database.Db.Debug().Where(models.Auction{UserID : user.ID}).Find(&Auctions).Error
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Fetch FAILED"})	
		return
	}
	c.JSON(http.StatusOK, gin.H{"Auctions": Auctions});
}
