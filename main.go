package main

import (
	
	"ipl/database"
	"github.com/gin-gonic/gin"
	"ipl/controllers"
	"ipl/middleware"
	"ipl/firebase"

)

func Initialize() {
	database.InitialiseDatabase()
	firebase.InitializeFirebase()
}


func main() {
	Initialize()
	router := gin.Default();

	admin := router.Group("/admin")
	{
		admin.POST("/signup",controllers.SignUpAdmin);
		admin.POST("/signin",controllers.SignInAdmin);
		
		admin.Use(middleware.AdminAuthMiddleware());
		admin.POST("/create-auction",controllers.CreateAuction);
		admin.POST("/get-all-auctions",controllers.AdminGetAllAuctions);
		admin.POST("/get-auction",controllers.AdminGetAuction);
		
	}

	userRoutes := router.Group("/user")
	{
		userRoutes.POST("/signup",controllers.SignUpUser);
		userRoutes.POST("/signin",controllers.SignInUser);
		
		userRoutes.Use(middleware.UserAuthMiddleware());
		
		userRoutes.POST("/create-bid",controllers.UserCreateBid);
		userRoutes.POST("/get-all-auctions",controllers.UserGetAllAuctions);
		userRoutes.POST("/get-auction",controllers.UserGetAuction);
		userRoutes.POST("/get-user-auctions",controllers.UserGetAllUserAuction)
	}
	router.Run(":3002");
}