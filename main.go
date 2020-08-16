package main

import (
	"github.com/gin-contrib/cors"
	"ipl/database"
	"github.com/gin-gonic/gin"
	"ipl/controllers"
	"ipl/middleware"
	"ipl/firebase"

	"time"
	"fmt"

)

func Initialize() {
	database.InitialiseDatabase()
	firebase.InitializeFirebase()
}


func main() {
	Initialize()
	router := gin.Default();

	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"PUT", "PATCH","POST","GET","DELETE"},
        AllowHeaders:     []string{"Origin","Content-Length","Content-Type","Authorization","AuthType"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

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

	go func (){
		for _ = range time.Tick(10 * time.Second) { 
			go database.BidWatchDogFucntion()
			fmt.Printf("Exec")
		 }
	}()


	router.Run(":3002");

	
	

}