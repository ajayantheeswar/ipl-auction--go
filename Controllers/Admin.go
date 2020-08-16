package controllers

import (
	"context"
	"fmt"
	"io"
	"ipl/database"
	"ipl/firebase"
	"ipl/models"
	"ipl/utils"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type POJOSignUP struct {
	Name     string
	Email    string
	Password string
	AuthType string
	Token    string `json:",omitempty"`
}

func SignUpAdmin(c *gin.Context) {
	var requestBody POJOSignUP
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Auth Status": "FAIL"})
		panic(err)
		return
	}

	var user = models.AdminUser{}

	if requestBody.AuthType == "Google" {
		Payload, err := utils.VerifyIDToken(requestBody.Token)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Auth Status": "FAIL"})
			panic(err)
			return
		}
		user.Email = Payload.Claims["email"].(string)
		user.Name = Payload.Claims["name"].(string)
		user.AuthType = "Google"
	} else {
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
		c.JSON(http.StatusConflict, gin.H{"Auth Status": "FAIL", "error": "The Account Already Exists"})
		panic(err)
		return
	} else {
		// User Not Exists
		token, err := utils.CreateToken(fmt.Sprint(user.ID))
		if err != nil {
			panic(err)
			c.JSON(http.StatusBadRequest, gin.H{"Auth Status": "FAIL"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Token": token,
			"Name": user.Name, "Email": user.Email,
		})
	}

}

type POJOSignIN struct {
	Email    string
	Password string
	AuthType string
	Token    string `json:",omitempty"`
}

func SignInAdmin(c *gin.Context) {
	var requestBody POJOSignIN
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Auth Status": "FAIL"})
		panic(err)
		return
	}

	var user models.AdminUser

	if requestBody.AuthType == "Google" {
		Payload, err := utils.VerifyIDToken(requestBody.Token)

		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"Auth Status": "FAIL"})
			panic(err)
			return
		}
		requestBody.Email = Payload.Claims["email"].(string)
	}
	// Check if the user Already Present
	err = database.Db.Where(&models.AdminUser{Email: requestBody.Email}).First(&user).Error

	if gorm.IsRecordNotFoundError(err) {
		// User Does Not EXISTS ... RETURN
		c.JSON(http.StatusForbidden, gin.H{"Auth Status": "FAIL", "error": "Invalid Credientials"})
		panic(err)
		return
	} else {
		// User Exists
		if requestBody.AuthType != "Google" {
			if !utils.CheckPasswordHash(requestBody.Password, user.Password) {
				// Invalid Credientials
				c.JSON(http.StatusForbidden, gin.H{"Auth Status": "FAIL", "error": "Invalid Credientials"})
				panic(err)
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

type POJOCreateAuction struct {
	Name         string `form:"Name" binding:"required"`
	BattingStyle string `form:"BattingStyle" binding:"required"`
	Average      string `form:"Average" binding:"required"`
	Role         string `form:"Role" binding:"required"`
	Country      string `form:"Country" binding:"required"`
	Start        uint64 `form:"Start" binding:"required"`
	End          uint64 `form:"End" binding:"required"`
}

// CereateAuction - USES FORM DATA
func CreateAuction(c *gin.Context) {

	user := c.Request.Context().Value("UserId").(models.AdminUser)

	var form POJOCreateAuction
	var auction models.Auction

	file, err := c.FormFile("Profile")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "File FAILED"})
		panic(err)
		return
	}

	err = c.ShouldBind(&form)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Form FAILED"})
		panic(err)
		return
	}

	//Upload The File
	url, err := UploadFile(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Upload FAILED"})
		panic(err)
		return
	}
	getAuctionFromForm(&auction, &form)
	auction.Profile = url
	auction.AdminUserID = user.ID

	// Update the Database
	err = database.Db.Debug().Create(&auction).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Upload FAILED"})
		panic(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "Auction Created Successfully !"})
}

func UploadFile(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()

	if err != nil {
		return "", fmt.Errorf("File Error -", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*50)
	// Upload an object with storage.Writer.
	wc := firebase.Bucket.Object(file.Filename).NewWriter(ctx)
	if _, err := io.Copy(wc, f); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	result := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%v/o/%v?alt=media", firebase.Attrs.Name, file.Filename)

	return result, nil
}

func getAuctionFromForm(Auction *models.Auction, form *POJOCreateAuction) {
	Auction.Name = form.Name
	Auction.BattingStyle = form.BattingStyle
	Auction.Average = form.Average
	Auction.Role = form.Role
	Auction.Start = form.Start
	Auction.End = form.End
}

//

func AdminGetAllAuctions(c *gin.Context) {
	user := c.Request.Context().Value("UserId").(models.AdminUser)

	var Auctions []models.Auction

	err := database.Db.Debug().Where(models.Auction{AdminUserID: user.ID}).Find(&Auctions).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Fetch FAILED"})
		panic(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"Auctions": Auctions})

}

type POJOAdminGetAuction struct {
	AuctionID uint
}

func AdminGetAuction(c *gin.Context) {

	var requestBody POJOAdminGetAuction
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Request FAILED"})
		panic(err)
		return
	}
	var Auction models.Auction

	err = database.Db.Debug().Preload("Bids").First(&Auction, requestBody.AuctionID).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": "Fetch FAILED"})
		panic(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"Auction": Auction})

}
