package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mgsquare/go-ecommerce/controllers"
	"github.com/mgsquare/go-ecommerce/database"
	"github.com/mgsquare/go-ecommerce/middleware"
	"github.com/mgsquare/go-ecommerce/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	app := controllers.NewApplication(database.ProductData(database.Client, "Product"), database.UserData(database.Client, "User"))
	fmt.Println("Server started on port: ", port)
	router := gin.New()

	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/listcart", controllers.GetItemFromCart())
	router.POST("/addaddress", controllers.AddAddress())
	router.PUT("/edithomeaddress", controllers.EditHomeAddress())
	router.PUT("/editworkaddress", controllers.EditWorkAddress())
	router.GET("/deleteaddresses", controllers.DeleteAddress())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())
	log.Fatal(router.Run(":" + port))

}
