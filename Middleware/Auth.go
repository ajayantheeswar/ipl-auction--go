package middleware

import (
	"context"
	"errors"
	"ipl/database"
	"ipl/models"
	"ipl/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// AdminAuthMiddleware - To Auth Middleware
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		AuthHeader, err := GetAuthHeader(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		AuthType, err := GetAuthTypeHeader(c.Request)
		user := models.AdminUser{}
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		// Have both AuthType and AuthHeader

		if AuthType == "Google" {
			// Google SignIn Verification
			Payload,err := utils.VerifyIDToken(AuthHeader)
		
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"Auth Status": "FAIL"})
				return
			}
			var Email string = Payload.Claims["email"].(string)
			err = database.Db.Where(&models.AdminUser{Email: Email}).First(&user).Error
			if gorm.IsRecordNotFoundError(err) {
				// The User is Not Valid
				c.JSON(http.StatusUnauthorized, err.Error())
				c.Abort()
				return
			}
			ctx := context.WithValue(c.Request.Context(), "UserId", user)
			c.Request = c.Request.WithContext(ctx)
			c.Next()


		} else {
			userId, err := utils.GetAuthClaims(AuthHeader)
			if err != nil {
				c.JSON(http.StatusUnauthorized, err.Error())
				c.Abort()
				return
			}
			// Got the ADMIN USER ID
			userIntegerID, err := strconv.Atoi(userId)

			if err != nil {
				c.JSON(http.StatusUnauthorized, err.Error())
				c.Abort()
				return
			}
			err = database.Db.First(&user, userIntegerID).Error
			if gorm.IsRecordNotFoundError(err) {
				// The User is Not Valid
				c.JSON(http.StatusUnauthorized, err.Error())
				c.Abort()
				return
			}
			ctx := context.WithValue(c.Request.Context(), "UserId", user)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		
		}
		return
	}
}

func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		AuthHeader, err := GetAuthHeader(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		AuthType, err := GetAuthTypeHeader(c.Request)

		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		// Have both AuthType and AuthHeader
		user := models.User{}

		if AuthType == "Google" {
			Payload,err := utils.VerifyIDToken(AuthHeader)
		
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"Auth Status": "FAIL"})
				return
			}
			var Email string = Payload.Claims["email"].(string)
			err = database.Db.Where(&models.User{Email: Email}).First(&user).Error
			if gorm.IsRecordNotFoundError(err) {
				// The User is Not Valid
				c.JSON(http.StatusUnauthorized, err.Error())
				c.Abort()
				return
			}
			ctx := context.WithValue(c.Request.Context(), "UserId", user)
			c.Request = c.Request.WithContext(ctx)
			c.Next()

		} else {
			userId, err := utils.GetAuthClaims(AuthHeader)
			if err != nil {
				c.JSON(http.StatusUnauthorized, err.Error())
				c.Abort()
				return
			}
			// Got the ADMIN USER ID
			userIntegerID, err := strconv.Atoi(userId)

			if err != nil {
				c.JSON(http.StatusUnauthorized, err.Error())
				c.Abort()
				return
			}

			
			err = database.Db.First(&user, userIntegerID).Error
			if gorm.IsRecordNotFoundError(err) {
				// The User is Not Valid
				c.JSON(http.StatusUnauthorized, err.Error())
				c.Abort()
				return
			}
			ctx := context.WithValue(c.Request.Context(), "UserId", user)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		}
		return
	}
}

func GetAuthHeader(req *http.Request) (string, error) {
	bearToken := req.Header.Get("Authorization")

	if bearToken == "" {
		return "", errors.New("Authentication Header Not Found")
	}
	return bearToken, nil

}

func GetAuthTypeHeader(req *http.Request) (string, error) {
	authType := req.Header.Get("authType")

	if authType == "" {
		return "", errors.New("authType Header Not Found")
	}
	return authType, nil

}
